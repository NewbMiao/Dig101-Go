package main

import (
	"fmt"
)

// https://github.com/kat-co/concurrency-in-go-src/blob/4e55fd7f3f5b9c5efc45a841702393a1485ba206/concurrency-patterns-in-go/the-bridge-channel/fig-bridge-channel.go
var bridge = func(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range stream {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

var genVals = func() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	done := make(chan interface{})
	defer close(done)
	for v := range bridge(done, genVals()) {
		fmt.Printf("%v ", v)
	}
}
