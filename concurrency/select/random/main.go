package main

import (
	"fmt"
)

// https://github.com/kat-co/concurrency-in-go-src/blob/master/gos-concurrency-building-blocks/the-select-statement/fig-select-uniform-distribution.go
func main() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}
