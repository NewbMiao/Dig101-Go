package main

import (
	"testing"
)

func repeat(
	done <-chan interface{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func BenchmarkOr(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	num := 100
	streams := make([]<-chan interface{}, num)
	for i := range streams {
		streams[i] = repeat(done, []int{1, 2, 3})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-Or(streams...)
	}
}

func BenchmarkOrInReflect(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	num := 100
	streams := make([]<-chan interface{}, num)
	for i := range streams {
		streams[i] = repeat(done, []int{1, 2, 3})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-OrInReflect(streams...)
	}
}
