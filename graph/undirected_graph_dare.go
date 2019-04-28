package graph

// Checking for cycles in a graph

// Cycle interface
type Cycle interface {
	HasCycle() bool
}

type cycle struct {
	marked   []bool
	hasCycle bool
}

// NewCycle creates cycle instance
func NewCycle(g Graph) Cycle {
	c := &cycle{marked: make([]bool, g.V())}
	for s := 0; s < g.V(); s++ {
		if !c.marked[s] {
			c.dfs(g, s, s)
		}
	}
	return c
}

func (c *cycle) dfs(g Graph, v, u int) {
	c.marked[v] = true
	for _, w := range g.Adj(v) {
		if !c.marked[w] {
			c.dfs(g, w, v)
		} else if w != u {
			c.hasCycle = true
		}
	}
}

func (c cycle) HasCycle() bool {
	return c.hasCycle
}

// BiPartite graph

// BiPartite interface
type BiPartite interface {
	IsBipartite() bool
}

type bipartite struct {
	marked      []bool
	color       []bool
	isBipartite bool
}

// NewBipartite creates a BiPartite instance
func NewBipartite(g Graph) BiPartite {
	b := &bipartite{
		marked:      make([]bool, g.V()),
		color:       make([]bool, g.V()),
		isBipartite: true,
	}
	for s := 0; s < g.V(); s++ {
		if !b.marked[s] {
			b.dfs(g, s)
		}
	}

	return b
}

func (b *bipartite) dfs(g Graph, v int) {
	b.marked[v] = true
	for _, w := range g.Adj(v) {
		if !b.marked[w] {
			b.color[w] = !b.color[v]
			b.dfs(g, w)
		} else if b.color[w] == b.color[v] {
			b.isBipartite = false
		}
	}
}

func (b bipartite) IsBipartite() bool {
	return b.isBipartite
}

// TODO Eulerian path

// EulerianCycle interface
type EulerianCycle interface {
	EulerianPath() []int
}

type eulerian struct{}

// NewEulerianPath creates an eulerian instance
func NewEulerianPath(g Graph) EulerianCycle {
	e := &eulerian{}
	return e
}

func (e *eulerian) dfs(g Graph, v int) {
	// TODO
}

func (e eulerian) EulerianPath() []int {
	// TODO
	return nil
}
