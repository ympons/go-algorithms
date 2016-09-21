package unionfind

import (
	"fmt"
	"io"
)

// UnionFind Api - Algorithm for Dynamic Connectivity
type (
	UnionFind interface {
		// add connection between two objects
		Union(int, int)
		// check if two objects are connected
		Connected(int, int) bool
		// return component identifier for the object (0 to N-1)
		Find(int) int
		// number of components
		Count() int
	}
	// NewUnionFind creates a UnionFind instance
	NewUnionFind func(n int) UnionFind
)

// NewUnionFindFromReader creates a new UF from a reader
func NewUnionFindFromReader(r io.Reader, create NewUnionFind) UnionFind {
	var n int

	fmt.Fscanf(r, "%d", &n)

	uf := create(n)
	for {
		var p, q int
		_, err := fmt.Fscanf(r, "%d %d", &p, &q)
		if err != nil {
			if err != io.EOF {
				return nil
			}
			break
		}
		if !uf.Connected(p, q) {
			uf.Union(p, q)
		}
	}

	return uf
}

// Quick-Find [Eager approach]
// Drawbacks:
// 	- Union too expensive (N array accesses)
// 	- Trees are flat, but to expensive to keep them flat
type qfUF struct {
	id    []int
	count int
}

// NewQuickFindUF creates a qfUF UnionFind instance
func NewQuickFindUF(n int) UnionFind {
	id := make([]int, n)
	// set id of each object to itself (N array accesses)
	for i := 0; i < n; i++ {
		id[i] = i
	}
	return &qfUF{id, n}
}

func (qf qfUF) Connected(p, q int) bool {
	// check whether p and q are in the same component (2 array accesses)
	return qf.id[p] == qf.id[q]
}

func (qf *qfUF) Union(p, q int) {
	pid, qid := qf.id[p], qf.id[q]
	if pid == qid {
		return
	}
	// change all entries with id[p] to id[q] (at most 2N+2 array accesses)
	for i := 0; i < len(qf.id); i++ {
		if qf.id[i] == pid {
			qf.id[i] = qid
		}
	}
	qf.count--
}

func (qf qfUF) Find(p int) int {
	return qf.id[p]
}

func (qf qfUF) Count() int {
	return qf.count
}

// Quick-Union [Lazy approach]
// Interpretation: id[i] is parent of i
// 	Root of i is id[id[id[...id[i]...]]]
// Drawbacks:
// 	- Trees can get tall
// 	- Find too expensive (could be N array accesses)
type quUF struct {
	id    []int
	count int
}

// NewQuickUnionUF creates a quUF UnionFind instance
func NewQuickUnionUF(n int) UnionFind {
	id := make([]int, n)
	// set id of each object to itself (N array accesses)
	for i := 0; i < n; i++ {
		id[i] = i
	}
	return &quUF{id, n}
}

func (qu quUF) Connected(p, q int) bool {
	// check if p and q have the same root (depth of p and q array accesses)
	return qu.Find(p) == qu.Find(q)
}

func (qu *quUF) Find(p int) int {
	// chase parent pointers until reach root (depth of i array accesses)
	for p != qu.id[p] {
		// this extra line is for path compression by halving
		// Keeps tree almost completely flat
		qu.id[p] = qu.id[qu.id[p]]

		p = qu.id[p]
	}
	return p
}

func (qu *quUF) Union(p, q int) {
	i, j := qu.Find(p), qu.Find(q)
	if i == j {
		return
	}
	// change root of p to point to root of q (depth of p and q array accesses)
	qu.id[i] = j
	qu.count--
}

func (qu quUF) Count() int {
	return qu.count
}

// Weighted Quick-Union
// 	- Modify Quick-Union to avoid tall trees
//	- Keep track of size of each tree (number of objects)
// 	- Balance by linking root of smaller tree to root of larger tree
type wquUF struct {
	quUF
	sz []int
}

// NewWeightedQuickUnionUF creates a wquUF UnionFind instance
func NewWeightedQuickUnionUF(n int) UnionFind {
	qu := NewQuickUnionUF(n).(*quUF)
	sz := make([]int, n)
	return &wquUF{*qu, sz}
}

func (wqu *wquUF) Union(p, q int) {
	i, j := wqu.Find(p), wqu.Find(q)
	if i == j {
		return
	}

	// make root of smaller size (sz) point to root of larger size
	if wqu.sz[i] < wqu.sz[j] {
		wqu.id[i] = j
	} else if wqu.sz[i] > wqu.sz[j] {
		wqu.id[j] = i
	} else {
		wqu.id[j] = i
		wqu.sz[i]++
	}
	wqu.count--
}
