package main

import (
	"fmt"
	"testing"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

var funcs = []struct {
	name string
	f    func(...<-chan interface{}) <-chan interface{}
}{
	{"reflection", fanInReflect},
	{"recursion", fanInRecur},
	{"goroutine", fanIn},
}

func TestFanIn(t *testing.T) {
	streams := make([]<-chan interface{}, 100)
	for _, fn := range funcs {
		t.Run(fn.name, func(t *testing.T) {
			done := make(chan struct{})
			defer close(done)
			for i := range streams {
				streams[i] = generator.AsStream(done, []int{1, 2, 3})
			}
			ch := fn.f(streams...)
			for range ch {
			}
		})
	}
}

func BenchmarkFanIn(b *testing.B) {
	for _, fn := range funcs {
		for n := 1; n <= 1024; n *= 2 {
			b.Run(fmt.Sprintf("%s/%d", fn.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					done := make(chan struct{})
					defer close(done)
					streams := make([]<-chan interface{}, n)
					for i := range streams {
						streams[i] = generator.AsStream(done, []int{1, 2, 3})
					}
					b.StartTimer()
					ch := fn.f(streams...)
					for range ch {
					}
				}
			})
		}
	}
}
