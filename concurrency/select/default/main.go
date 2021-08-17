package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		start := time.Now()
		select {
		case <-done:
			break loop
		default:
			// Simulate work
			fmt.Printf("In default after %v\n\n", time.Since(start))
			workCounter++
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("Achieved %v cycles of work before signaled to stop.\n", workCounter)
}
