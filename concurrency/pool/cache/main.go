package main

import (
	"sync"
)

// https://github.com/kat-co/concurrency-in-go-src/blob/master/gos-concurrency-building-blocks/the-sync-package/pool/fig-sync-pool.go
const maxSize = 4096

func main() {
	allocCnt := 0
	cacheCnt := 4
	p := &sync.Pool{
		New: func() interface{} {
			allocCnt++
			b := make([]byte, maxSize)
			return &b
		},
	}
	for i := 0; i < cacheCnt; i++ {
		p.Put(p.New())
	}

	numGoroutines := 1024
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			m, _ := p.Get().(*[]byte)
			defer p.Put(m)
		}()
	}
	wg.Wait()
	println("alloc memory times: ", allocCnt)
}
