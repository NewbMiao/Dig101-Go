package main

import (
	"fmt"
	"testing"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

func TestBridge(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	streams := generator.GenerateChanStream(10)
	for range Bridge(done, streams) {
	}
}

func BenchmarkBridge(b *testing.B) {
	for n := 1; n <= 1024; n *= 2 {
		b.Run(fmt.Sprintf("/%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				done := make(chan struct{})
				defer close(done)
				streams := generator.GenerateChanStream(i)
				b.StartTimer()
				for range Bridge(done, streams) {
				}
			}
		})
	}
}
