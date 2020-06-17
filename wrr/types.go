package wrr

import "sync"

type (
	// WRR describe the basic methods that used to be used for an different implementations of a Weighted Round-Robin algorithm.
	// The most common use of any instance that implements it — request balancers.
	WRR interface {
		// Add adds a new entry (i.e. a connection) to the weighted set that implements Weighted RR algorithm.
		Add(item interface{}, weight int64)
		// Next pops the next item returned by WRR algorithm from the weighted set.
		Next() interface{}
	}

	// edf

	edfEntry struct {
		item     interface{}
		idx      int
		weight   int64
		deadline float64
	}

	pqueue []*edfEntry

	edfRr struct {
		mu        sync.Mutex
		q         pqueue
		itemsPool map[interface{}]*edfEntry
	}

	// edf

	// random

	weightedItem struct {
		item   interface{}
		weight int64
	}

	randomRr struct {
		mu           sync.RWMutex
		items        []*weightedItem
		sumOfWeights int64
	}

	// random
)
