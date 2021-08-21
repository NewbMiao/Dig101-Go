package main

import (
	"fmt"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

// https://github.com/kat-co/concurrency-in-go-src/blob/bb10a902ef1bcaf788d2c3ab9475ceb24f05c5fe/concurrency-patterns-in-go/the-tee-channel/fig-tee-channel.go

func tee(
	done <-chan struct{},
	in <-chan interface{},
) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range in {
			out1copy, out2copy := out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1copy <- val:
					out1copy = nil // 阻塞
				case out2copy <- val:
					out2copy = nil // 阻塞
				}
			}
		}
	}()
	return out1, out2
}

func main() {
	done := make(chan struct{})
	defer close(done)

	out1, out2 := tee(done, generator.TakeN(done, generator.Repeat(done, 1, 2, 3, 4), 4))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}
