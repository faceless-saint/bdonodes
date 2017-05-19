package types

import (
	"container/heap"
)

type NodeHeap []*Node

// Len returns the length of the NodeHeap.
func (m NodeHeap) Len() int {
	return len(m)
}

// Less returns true if element i is less than element j.
func (m NodeHeap) Less(i, j int) bool {
	return m[i].priority() < m[j].priority()
}

// Swap elements i and j in the NodeHeap, updating each index.
func (h NodeHeap) Swap(i, j int) bool {
	h[i], h[j] = h[j], h[i]
	h[i].setIndex(i)
	h[j].setIndex(j)
}

// Push the given Node onto the NodeHeap.
func (h *NodeHeap) Push(x interface{}) {
	item := Node(x)
	item.setIndex(len(*h))
	*h = append(*h, item)
}

// Pop the top Node from the NodeHeap.
func (h *NodeHeap) Pop() interface{} {
	item := (*h)[len(*h)-1]
	item.setIndex(-1)
	*h = (*h)[0:len(*h)-1]
	return item
}
