package quadtree

import "container/heap"

type PriorityQueue []*treeNode

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority (color diff) so we use ">" here.
	return pq[i].diff > pq[j].diff
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*treeNode)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

// BuildPQ initializes the priority queue and returns it
func BuildPQ() PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return pq
}
