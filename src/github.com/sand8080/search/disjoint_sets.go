package search

import (
	"fmt"
	"errors"
)

type DisjointSetInt struct {
	Parents map[int]int
	Weights map[int]int
}

func NewDisjointSetInt(size int) *DisjointSetInt {
	return &DisjointSetInt{
		Parents: make(map[int]int, size),
		Weights: make(map[int]int, size),
	}
}

func (ds *DisjointSetInt) getPath(id int) ([]int, int) {
	path := []int{id}
	root, ok := ds.Parents[id]
	if !ok {
		panic(fmt.Sprintf("Key %v not found in parents", id))
	}

	for {
		if root == path[len(path) - 1] {
			return path, root
		}
                path = append(path, root)
		root, ok = ds.Parents[root]
		if !ok {
			panic(fmt.Sprintf("Key %v not found in parents", id))
		}
	}
}

func (ds *DisjointSetInt) getOrCreateRoot(id int) int {
	if _, ok := ds.Parents[id]; !ok {
		ds.Parents[id] = id
		ds.Weights[id] = 1
		return id
	}

	path, root := ds.getPath(id)

	// Compress the path and return
	for ancestor := range path {
		ds.Parents[ancestor] = root
	}
	return root
}

func (ds * DisjointSetInt) getHeaviest(roots []int) (int, int, error) {
	if len(roots) == 0 {
		return 0, 0, errors.New("Can't find minimal value in the " +
			"empty collection")
	}

	max_idx := 0
	max, ok := ds.Weights[roots[max_idx]]
	if !ok {
		return 0, 0, errors.New(fmt.Sprintf("'%v' is not found in weights",
			roots[max_idx]))
	}

	for i := 1; i < len(roots); i++ {
		if max < roots[i] {
			max, max_idx = roots[i], i
		}
	}

	return max, max_idx, nil

}

func (ds *DisjointSetInt) getRoots(ids []int) []int {
	result := make([]int, len(ids))
	for idx, id := range ids {
		result[idx] = ds.getOrCreateRoot(id)
	}
	return result
}

func (ds *DisjointSetInt) Union(ids []int) {
	if len(ids) == 0 {
		return
	}
	roots := ds.getRoots(ids)
	fmt.Printf("Union for int is called: %v\n", roots)
}