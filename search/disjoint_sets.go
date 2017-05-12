package search

import (
	"errors"
	"fmt"
	"sync"
)

type DisjointSetInt struct {
	Parents map[int]int
	Weights map[int]int
	sync.Mutex
}

func NewDisjointSetInt(size int) *DisjointSetInt {
	return &DisjointSetInt{
		Parents: make(map[int]int, size),
		Weights: make(map[int]int, size),
	}
}

func (ds *DisjointSetInt) getPath(id int) ([]int, int, error) {
	path := []int{id}
	root, ok := ds.Parents[id]
	if !ok {
		msg := fmt.Sprintf("Key %v not found in parents", id)
		return nil, 0, errors.New(msg)
	}

	for {
		if root == path[len(path)-1] {
			return path, root, nil
		}
		path = append(path, root)
		root, ok = ds.Parents[root]
		if !ok {
			panic(fmt.Sprintf("Key %v not found in parents", id))
		}
	}
}

func (ds *DisjointSetInt) getOrCreateRoot(id int) (int, error) {
	if _, ok := ds.Parents[id]; !ok {
		ds.Parents[id] = id
		ds.Weights[id] = 1
		return id, nil
	}

	path, root, err := ds.getPath(id)
	if err != nil {
		return 0, err
	}

	// Compress the path and return
	for _, ancestor := range path {
		ds.Parents[ancestor] = root
	}
	return root, nil
}

func (ds *DisjointSetInt) getHeaviest(ids, roots []int) (int, int, error) {

	if len(roots) == 0 {
		return 0, 0, errors.New("Can't find minimal value in the " +
			"empty collection")
	}

	noWeightError := func(key int) (int, int, error) {
		err := errors.New(fmt.Sprintf("'%v' is not found in weights", key))
		return 0, 0, err
	}

	max_idx := 0
	max, ok := ds.Weights[roots[max_idx]]
	if !ok {
		return noWeightError(roots[max_idx])
	}

	for root_idx := 1; root_idx < len(roots); root_idx++ {
		root := roots[root_idx]
		root_weight, ok := ds.Weights[root]
		if !ok {
			return noWeightError(root_idx)
		}
		if max < root_weight ||
			(max == root_weight && ids[max_idx] < ids[root_idx]) {
			max, max_idx = root_weight, root_idx
		}
	}
	return roots[max_idx], max_idx, nil
}

func (ds *DisjointSetInt) getRoots(ids []int) ([]int, error) {
	result := make([]int, len(ids))
	for idx, id := range ids {
		root, err := ds.getOrCreateRoot(id)
		if err != nil {
			return nil, err
		}
		result[idx] = root
	}
	return result, nil
}

func (ds *DisjointSetInt) Union(ids []int) error {
	ds.Lock()
	defer ds.Unlock()

	if len(ids) == 0 {
		return nil
	}
	roots, err := ds.getRoots(ids)
	if err != nil {
		return err
	}

	heaviest, _, err := ds.getHeaviest(ids, roots)
	if err != nil {
		return err
	}

	// Increase weight of root obj only if it was't done yet
	for root_idx, root := range roots {
		id := ids[root_idx]
		if root != heaviest {
			_, effective_root, err := ds.getPath(id)
			if err != nil {
				return err
			}

			if effective_root != heaviest {
				root_weight := ds.Weights[root]
				ds.Weights[heaviest] += root_weight
				ds.Parents[root] = heaviest
			}
		}
	}
	return nil
}

// Type for emitting groups of ids
type group struct {
	ids    []int
	weight int
}

func (ds *DisjointSetInt) EmitGroups() <-chan []int {
	ch := make(chan []int)
	go func() {
		defer close(ch)
		ds.Lock()
		defer ds.Unlock()

		ready_to_dump := make(map[int]group)

		for id := range ds.Parents {
			_, root, err := ds.getPath(id)
			if err != nil {
				panic(err)
			}

			g, ok := ready_to_dump[root]
			if ok {
				g.ids = append(g.ids, id)
			} else {
				g.ids = []int{id}
				g.weight = ds.Weights[root]
			}

			g.weight -= 1
			ready_to_dump[root] = g

			// All children are collected
			if g.weight == 0 {
				ch <- g.ids
				delete(ready_to_dump, root)
			} else if g.weight < 0 {
				panic("Group emission failed: weight less than 0")
			}
		}
		if len(ready_to_dump) > 0 {
			panic(fmt.Sprintf("Not all groups dumped: %v", ready_to_dump))
		}
	}()
	return ch
}
