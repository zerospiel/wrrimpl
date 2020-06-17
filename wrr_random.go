package wrrimpl

import (
	"math/rand"
	"time"
)

// assert implementation
var _ WRR = (*randomRr)(nil)

var rsource = rand.New(rand.NewSource(time.Now().UnixNano()))

// NewRandom creates returns and instance of a struct that implements WRR logic.
// This constructor basics on random item's weight picking within all weights of the set of items.
func NewRandom() WRR {
	return &randomRr{}
}

func (r *randomRr) Add(item interface{}, weight int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rItem := &weightedItem{
		item:   item,
		weight: weight,
	}
	r.items = append(r.items, rItem)
	r.sumOfWeights += weight
}

func (r *randomRr) Next() interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.sumOfWeights == 0 {
		return nil
	}

	randomWeight := rsource.Int63n(r.sumOfWeights)
	for _, v := range r.items {
		randomWeight = randomWeight - v.weight
		if randomWeight < 0 {
			return v.item
		}
	}

	return r.items[len(r.items)-1].item
}
