package quadtree

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue(t *testing.T) {

	diffs := []float64{0, 10, 3, 5}
	pq := BuildPQ()

	nodes := make([]*treeNode, len(diffs))
	for i := range nodes {
		nodes[i] = new(treeNode)
		nodes[i].diff = diffs[i]
	}

	for _, node := range nodes {
		heap.Push(&pq, node)
	}

	indices := []int{1, 3, 2, 0}
	idx := 0

	for pq.Len() > 0 {
		if node := heap.Pop(&pq).(*treeNode); node.diff != diffs[indices[idx]] {
			t.Errorf("diff mismatch ! %f ; %f", node.diff, diffs[indices[idx]])
		}

		idx++
	}

}
