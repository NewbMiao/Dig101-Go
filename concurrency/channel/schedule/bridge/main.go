package main

import (
	"fmt"
	"runtime"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

// https://github.com/kat-co/concurrency-in-go-src/blob/4e55fd7f3f5b9c5efc45a841702393a1485ba206/concurrency-patterns-in-go/the-bridge-channel/fig-bridge-channel.go
func Bridge(
	done <-chan struct{},
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
			for val := range generator.OrDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func main() {
	fmt.Println("任务进行中，当前协程数:", runtime.NumGoroutine())

	done := make(chan struct{})
	defer close(done)
	streams := generator.GenerateChanStream(10)
	for v := range Bridge(done, streams) {
		fmt.Printf("%v ", v)
	}

	fmt.Println("\n任务结束，当前协程数:", runtime.NumGoroutine())
}
