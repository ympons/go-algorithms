package unionfind

import (
	"os"
	"testing"
)

type Fatalistic interface {
	Fatal(...interface{})
}

func unionFindTest(t Fatalistic, create NewUnionFind) UnionFind {
	f, err := os.Open("../data/tinyUF.txt")
	if err != nil {
		t.Fatal(err)
	}
	uf := NewUnionFindFromReader(f, create)
	if uf == nil {
		t.Fatal("Error creating UnionFind from File")
	}
	return uf
}

func TestNewQuickFindUF(t *testing.T) {
	qf := unionFindTest(t, NewQuickFindUF)
	t.Log(qf)
}

func TestNewQuickUnionUF(t *testing.T) {
	qu := unionFindTest(t, NewQuickUnionUF)
	t.Log(qu)
}

func TestNewWeightedQuickUnionUF(t *testing.T) {
	wqu := unionFindTest(t, NewWeightedQuickUnionUF)
	t.Log(wqu)
}

func TestAllUnionFind(t *testing.T) {
	cqf := unionFindTest(t, NewQuickFindUF).Count()
	cqu := unionFindTest(t, NewQuickUnionUF).Count()
	cwqu := unionFindTest(t, NewWeightedQuickUnionUF).Count()

	if cqf != cqu || cqf != cwqu {
		t.Error("All the UnionFind variants should return the same number of components")
	}
}
