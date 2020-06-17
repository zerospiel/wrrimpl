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

// TODO: 
```

[1]: https://en.wikipedia.org/wiki/Earliest_deadline_first_scheduling
[2]: https://en.wikipedia.org/wiki/Weighted_round_robin
[3]: https://pkg.go.dev/google.golang.org/grpc@v1.29.1/balancer/base?tab=doc