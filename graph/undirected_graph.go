package graph

import (
	"fmt"
	"io"

	"github.com/ympons/go-algorithms/queue"
	"github.com/ympons/go-algorithms/stack"
)

// Undirected Graph API

// Graph interface
type Graph interface {
	// number of vertices
	V() int
	// number of edges
	E() int
	// add an edge
	AddEdge(int, int)
	// vertices adjacent to a vertex
	Adj(int) []int
	// string representation
	String() string
}

// Factory represents a Graph factory interface
type Factory interface {
	// Create an empty graph with n vertices
	NewGraph(int) Graph
	// Create a graph from a reader
	NewGraphFromReader(io.Reader) Graph
}

type graph struct {
	v, e int
	adj  [][]int
}

// NewGraph creates a new graph
func NewGraph(n int) Graph {
	adj := make([][]int, n)
	for v := 0; v < n; v++ {
		adj[v] = make([]int, 0, n)
	}
	return &graph{v: n, e: 0, adj: adj}
}

// NewGraphFromReader creates a graph from a reader
func NewGraphFromReader(r io.Reader) Graph {
	var n, E int

	fmt.Fscanf(r, "%d", &n)
	fmt.Fscanf(r, "%d", &E)

	g := NewGraph(n)
	for i := 0; i < E; i++ {
		var v, w int
		fmt.Fscanf(r, "%d %d", &v, &w)
		g.AddEdge(v, w)
	}
	return g
}

// Add an edge v-w
func (g *graph) AddEdge(v, w int) {
	g.adj[v] = append(g.adj[v], w)
	g.adj[w] = append(g.adj[w], v)
	g.e++
}

// Number of vertices
func (g graph) V() int {
	return g.v
}

// Number of edges
func (g graph) E() int {
	return g.e
}

// Vertices adjacent to v
func (g graph) Adj(v int) []int {
	return g.adj[v]
}

// String representation
func (g graph) String() string {
	return fmt.Sprintf("graph: V: %d E: %d -> %+v", g.V(), g.E(), g.adj)
}

// Graph-processing:
// compute the degree of v
func degree(G Graph, v int) int {
	degree := 0
	for _ = range G.Adj(v) {
		degree++
	}
	return degree
}

// compute maximum degree
func maxDegree(G Graph) int {
	max := 0
	for v := 0; v < G.V(); v++ {
		if d := degree(G, v); d > max {
			max = d
		}
	}
	return max
}

// compute average degree
func avgDegree(G Graph) int {
	return 2 * G.E() / G.V()
}

// count self-loops
func numberOfSelfLoops(G Graph) int {
	count := 0
	for v := 0; v < G.V(); v++ {
		for _, w := range G.Adj(v) {
			if v == w {
				count++
			}
		}
	}
	return count / 2
}

// Graph processing

// Paths represents a Graph Paths interface
type Paths interface {
	HasPathTo(int) bool
	PathTo(int) []int
}

// Depth-first search
type dfsPaths struct {
	marked []bool
	edgeTo []int
	s      int
}

// NewDFSPaths creates a DFS path
func NewDFSPaths(g Graph, s int) Paths {
	paths := &dfsPaths{
		marked: make([]bool, g.V()),
		edgeTo: make([]int, g.V()),
		s:      s,
	}
	paths.dfs(g, s)

	return paths
}

func (p *dfsPaths) dfs(g Graph, v int) {
	p.marked[v] = true
	for _, w := range g.Adj(v) {
		if !p.marked[w] {
			p.dfs(g, w)
			p.edgeTo[w] = v
		}
	}
}

func (p *dfsPaths) iterativeDfs(g Graph, v int) {
	p.marked[v] = true
	s := stack.NewStack()
	for _, w := range g.Adj(v) {
		s.Push(w)
	}
	for !s.IsEmpty() {
		w, _ := s.Pop()
		for _, u := range g.Adj(w) {
			if !p.marked[u] {
				p.marked[u] = true
				s.Push(u)
			}
		}
	}
}

func (p *dfsPaths) HasPathTo(v int) bool {
	return p.marked[v]
}

func (p *dfsPaths) PathTo(v int) []int {
	if !p.marked[v] {
		return nil
	}

	path := make([]int, 0, len(p.marked))
	for x := v; x != p.s; x = p.edgeTo[x] {
		path = append(path, x)
	}
	path = append(path, p.s)

	return path
}

// Breadth-first search
type bfsPaths struct {
	marked []bool
	edgeTo []int
	s      int
}

// NewBFSPaths creates a BFS path
func NewBFSPaths(g Graph, s int) Paths {
	paths := &bfsPaths{
		marked: make([]bool, g.V()),
		edgeTo: make([]int, g.V()),
		s:      s,
	}
	paths.bfs(g, s)

	return paths
}

func (p *bfsPaths) bfs(g Graph, s int) {
	p.marked[s] = true
	q := queue.NewQueue()
	q.Enqueue(s)
	for !q.IsEmpty() {
		v, _ := q.Dequeue()
		for _, w := range g.Adj(v) {
			if !p.marked[w] {
				q.Enqueue(w)
				p.marked[w] = true
				p.edgeTo[w] = v
			}
		}
	}
}

func (p *bfsPaths) HasPathTo(v int) bool {
	return p.marked[v]
}

func (p *bfsPaths) PathTo(v int) []int {
	if !p.marked[v] {
		return nil
	}

	path := make([]int, 0, len(p.marked))
	for x := v; x != p.s; x = p.edgeTo[x] {
		path = append(path, x)
	}
	path = append(path, p.s)

	return path
}

// Connected Components

// CC represents a Connected Components interface
type CC interface {
	Connected(v, w int) bool
	Count() int
	ID(v int) int
}

type cc struct {
	marked []bool
	id     []int
	count  int
}

// NewCC creates a Connected Component instance
func NewCC(g Graph) CC {
	cc := &cc{
		marked: make([]bool, g.V()),
		id:     make([]int, g.V()),
		count:  0,
	}
	for v := 0; v < g.V(); v++ {
		if !cc.marked[v] {
			cc.dfs(g, v)
			cc.count++
		}
	}

	return cc
}

func (c *cc) dfs(g Graph, v int) {
	c.marked[v] = true
	c.id[v] = c.count
	for _, w := range g.Adj(v) {
		if !c.marked[w] {
			c.dfs(g, w)
		}
	}
}

func (c cc) Connected(v, w int) bool {
	return c.id[v] == c.id[w]
}

func (c cc) Count() int {
	return c.count
}

func (c cc) ID(v int) int {
	return c.id[v]
}
