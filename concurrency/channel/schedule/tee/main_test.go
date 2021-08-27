package main

import (
	"fmt"
	"testing"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

func TestBridge(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	streams := generator.AsStream(done, []interface{}{1, 2, 3}...)
	out1, out2 := Tee(done, streams)
	for v := range out1 {
		if v != <-out2 {
			panic("test failed")
		}
	}
}

func BenchmarkBridge(b *testing.B) {
	for n := 1; n <= 1024; n *= 2 {
		b.Run(fmt.Sprintf("Tee-%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				done := make(chan struct{})
				defer close(done)
				streams := generator.AsStream(done, []interface{}{1, 2, 3}...)
				b.StartTimer()
				out1, out2 := Tee(done, streams)

				for v := range out1 {
					if v != <-out2 {
						panic("Not matched")
					}
				}
			}
		})
	}
}
