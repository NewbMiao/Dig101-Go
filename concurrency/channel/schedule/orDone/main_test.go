package main

import (
	"testing"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

var num = 500

func BenchmarkOr(b *testing.B) {
	done := make(chan struct{})
	defer close(done)
	streams := make([]<-chan interface{}, num)
	for i := range streams {
		streams[i] = generator.Repeat(done, []int{1, 2, 3})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-Or(streams...)
	}
}

func BenchmarkOrInReflect(b *testing.B) {
	done := make(chan struct{})
	defer close(done)
	streams := make([]<-chan interface{}, num)
	for i := range streams {
		streams[i] = generator.Repeat(done, []int{1, 2, 3})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-OrInReflect(streams...)
	}
}
