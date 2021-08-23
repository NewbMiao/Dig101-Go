package main

import (
	"fmt"
	"testing"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

var funcs = []struct {
	name string
	f    func(<-chan interface{}, []chan interface{})
}{
	{"reflection", fanOutReflect},
	{"iteration", fanOut},
	{"iterationAsync", fanOutAsync},
}

func TestFanOut(t *testing.T) {
	for _, fn := range funcs {
		t.Run(fn.name, func(t *testing.T) {
			done := make(chan struct{})
			defer close(done)
			in := generator.AsStream(done, []int{1, 2, 3})
			streams := make([]chan interface{}, 100)
			for i := range streams {
				streams[i] = make(chan interface{})
			}
			fn.f(in, streams)
			for _, v := range streams {
				t.Log(v)
			}
		})
	}
}

func BenchmarkFanOut(b *testing.B) {
	for _, fn := range funcs {
		for n := 1; n <= 1024; n *= 2 {
			b.Run(fmt.Sprintf("%s/%d", fn.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					done := make(chan struct{})
					defer close(done)
					in := generator.AsStream(done, []int{1, 2, 3})
					streams := make([]chan interface{}, n)
					for i := range streams {
						streams[i] = make(chan interface{})
					}
					b.StartTimer()
					fn.f(in, streams)
					for i := range streams {
						<-streams[i]
					}
				}
			})
		}
	}
}
