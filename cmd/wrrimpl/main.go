package main

import (
	"fmt"

	"github.com/zerospiel/wrrimpl/wrr"
)

type foobar struct {
	s string
}

func main() {
	edf := wrr.NewEDF()
	edf.Add(foobar{s: "hello"}, 70)
	edf.Add(foobar{s: "world"}, 30)

	random := wrr.NewRandom()
	random.Add(foobar{s: "hello"}, 70)
	random.Add(foobar{s: "world"}, 30)

	var (
		hedf, wedf, hrand, wrand int
	)

	for i := 0; i < 10000; i++ {
		we := edf.Next().(foobar).s
		wr := random.Next().(foobar).s
		if we == "hello" {
			hedf++
		} else {
			wedf++
		}
		if wr == "hello" {
			hrand++
		} else {
			wrand++
		}
	}

	fmt.Printf("%d queries\nEDF: `hello`: %d; `world`: %d\nRandom: `hello`: %d; `world`: %d\n", 10000, hedf, wedf, hrand, wrand)
}
