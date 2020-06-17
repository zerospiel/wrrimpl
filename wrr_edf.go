package wrrimpl

import (
	"container/heap"
)

// assert implementation
var _ WRR = (*edfRr)(nil)

// heap.Interface impl (min)

func (q pqueue) Len() int           { return len(q) }
func (q pqueue) Less(i, j int) bool { return q[i].deadline < q[j].deadline }
func (q pqueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].idx = i
	q[j].idx = j
}
func (q *pqueue) Push(x interface{}) {
	*q = append(*q, x.(*edfEntry))
}
func (q *pqueue) Pop() interface{} {
	old := *q
	*q = old[0 : len(old)-1]
	return old[len(old)-1]
}

// heap.Interface impl (min)

// NewEDF creates returns and instance of a struct that implements WRR logic.
// This constructor basics on EDF algorithm (https://en.wikipedia.org/wiki/Earliest_deadline_first_scheduling).
func NewEDF() WRR {
	return &edfRr{
		itemsPool: make(map[interface{}]*edfEntry),
	}
}

// currentDL returns current time to proceed in EDF.
// WARN: bypassing mutex by value due to unused and this is more optimal memory mgmt operation.
//nolint:copylocks
func (e edfRr) currentDL() float64 {
	if len(e.q) == 0 {
		return 0.
	}
	return e.q[0].deadline
}

// Add adds an item to the priority queue that prioritize items closes to its deadline id desc order.
// Similar to upsert operation behavior.
func (e *edfRr) Add(item interface{}, weight int64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// NOTE: not sure if this check is optimal due to mostly adding items
	// need general escape analyse (for currentDL()) + benchmarks
	// go tool compile -S
	v, ok := e.itemsPool[item]
	if ok {
		v.weight += weight
		return
	}
	v = &edfEntry{
		idx:      e.q.Len(),
		deadline: e.currentDL() + 1./float64(weight),
		item:     item,
		weight:   weight,
	}
	e.itemsPool[item] = v
	heap.Push(&e.q, v)
}

// Next returns the next (closest to deadline) item to process.
func (e *edfRr) Next() interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.q) == 0 {
		return nil
	}
	v := e.q[0]
	v.deadline = e.currentDL() + 1./float64(v.weight)
	heap.Fix(&e.q, 0)
	return v.item
}

// Remove removes a given item from the set.
func (e *edfRr) Remove(item interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if v, ok := e.itemsPool[item]; ok {
		heap.Remove(&e.q, v.idx)
	}
}
