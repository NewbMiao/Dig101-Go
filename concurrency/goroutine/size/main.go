package main

import (
	"fmt"
	"runtime"
	"sync"
)

// https://github.com/kat-co/concurrency-in-go-src/blob/master/gos-concurrency-building-blocks/goroutines/fig-goroutine-size.go
func main() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	// var c <-chan interface{}
	c := make(<-chan interface{})
	var wg sync.WaitGroup
	// never stop goroutine
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e6
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	mUsedInKB := float64(after-before) / 1024
	fmt.Printf("Memory used(%.f goroutines): per goroutine: %.2fKB; all: %.2fGB\n",
		numGoroutines,
		mUsedInKB/numGoroutines,
		mUsedInKB/1024/1024)
}
