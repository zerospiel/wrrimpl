# [EDF][1] and Random implementations for [weighted round-robin][2] scheduling

Used in implementation of a custom gRPC balancers and/or pickers.

Example of usage in custom picker using default [`google.golang.org/grpc/balancer/base`][3] balancer implementation:
```golang
package main

import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const BalancerName = "my-balancer"

type (
    pickerBldr struct{}
)

func init() {
    balancer.Register(newPBuilder(BalancerName))
}

func newPBuilder(name string) balancer.Builder {
    return base.NewBalancerBuilderV2(name, &pickerBldr{}, base.Config{HealthCheck: true})
}

func (p *pickerBldr) Build(info base.PickerBuildInfo) balancer.V2Picker {
	pool := newConnSet()

	for conn, connInfo := range info.ReadySCs {
		// any way you can receive weight for a new SC
		w := connInfo.Address.Attributes.Value("weight").(int64)
		pool.add(conn, w)
	}

	return &wrrPicker{
		p: pool,
	}
}

type connSet struct {
	mu  sync.Mutex
	wrr wrrimpl.WRR
}

func newConnSet() *connSet {
	return &connSet{wrr: wrrimpl.NewEDF()}
}

func (cs *connSet) add(sc balancer.SubConn, w int64) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.wrr.Add(sc, int64(w))
}

func (cs *connSet) pick() (balancer.PickResult, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// wrr on edf
	sc := cs.wrr.Next().(balancer.SubConn)
	if sc == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	return balancer.PickResult{SubConn: sc}, nil
}

type wrrPicker struct {
	p *connSet
}

// Pick returns the connection to use for this RPC and related information.
func (p *wrrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	return p.p.pick()
}
```

[1]: https://en.wikipedia.org/wiki/Earliest_deadline_first_scheduling
[2]: https://en.wikipedia.org/wiki/Weighted_round_robin
[3]: https://pkg.go.dev/google.golang.org/grpc@v1.29.1/balancer/base?tab=doc