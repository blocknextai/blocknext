package dag

// Item is an element of a PriorityQueue, identified by ID and ordered by Priority.
type Item struct {
	ID       string
	Priority int
	index    int
}

// PriorityQueue is a min-heap of Items, ordered by priority and then by ID. It
// implements heap.Interface.
type PriorityQueue []*Item

// Len returns the number of items in the queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less reports whether the item at index i should sort before the item at index j.
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Priority == pq[j].Priority {
		return pq[i].ID < pq[j].ID
	}
	return pq[i].Priority < pq[j].Priority
}

// Swap exchanges the items at indices i and j.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds x to the queue. It is intended to be called by heap.Push.
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the minimum item from the queue. It is intended to be
// called by heap.Pop.
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
