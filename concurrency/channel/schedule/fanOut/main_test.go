package main

import (
	"testing"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

var num = 100

func BenchmarkFanOut(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	in := generator.Repeat(done, []interface{}{1, 2, 3}...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		streams := make([]chan interface{}, num)
		for i := range streams {
			streams[i] = make(chan interface{})
		}
		fanOut(in, streams, false)
	}
}

func BenchmarkFanOutReflect(b *testing.B) {
	done := make(chan struct{})
	defer close(done)
	in := generator.Repeat(done, []interface{}{1, 2, 3}...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		streams := make([]chan interface{}, num)
		for i := range streams {
			streams[i] = make(chan interface{})
		}
		fanOutReflect(in, streams)
	}
}
