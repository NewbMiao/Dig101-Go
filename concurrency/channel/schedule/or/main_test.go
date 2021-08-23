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
	{"reflection", OrReflection},
	{"recursion", OrRecur},        // 二分递归
	{"recursion2", OrRecurSimple}, // 普通递归
	{"goroutine", Or},
}

func TestOr(t *testing.T) {
	for _, f := range funcs {
		streams := make([]<-chan interface{}, 10)
		t.Run(f.name, func(t *testing.T) {
			done := make(chan struct{})
			defer close(done)
			for i := range streams {
				streams[i] = generator.AsStream(done, []interface{}{1})
			}
			<-f.f(streams...)
		})
	}
}

// see https://docs.google.com/spreadsheets/d/11lVkxeSC8dRcTNxi4FubI-_Hls-4btCD13NAubXiFIY/edit?usp=sharin
func BenchmarkOr(b *testing.B) {
	for _, f := range funcs {
		for n := 8; n <= 1024; n *= 2 {
			b.Run(fmt.Sprintf("%s/%d", f.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					done := make(chan struct{})
					defer close(done)
					streams := make([]<-chan interface{}, n)
					for i := range streams {
						streams[i] = generator.AsStream(done, []interface{}{1})
					}
					b.StartTimer()
					<-f.f(streams...)
				}
			})
		}
	}
}
