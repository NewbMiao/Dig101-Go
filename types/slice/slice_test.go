package main

// go test . -bench  BenchmarkCutSlice -benchmem   -memprofile=mem.out
// go tool pprof -http=localhost:8080 mem.out

import "testing"

type bigStruct struct {
	id   int
	data [1024]int
}

const loopCnt = 1000

// nolint: gochecknoinits // example
func init() {
	discardLog()
}

func BenchmarkCutSlicePointer(b *testing.B) {
	a := make([]T, loopCnt)
	for i := 0; i < loopCnt; i++ {
		a[i] = &bigStruct{id: i}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cutSlicePointer(a, 1, 800)
	}
}

func BenchmarkCutSlice(b *testing.B) {
	a := make([]T, loopCnt)
	for i := 0; i < loopCnt; i++ {
		a[i] = &bigStruct{id: i}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cutSlice(a, 1, 800)
	}
}

func BenchmarkCopy(b *testing.B) {
	src := make(Slice, 10240)
	src[1024] = 1
	dst := make(Slice, 10240)
	for i := 0; i < b.N; i++ {
		_ = copy(src[:], src)
		_ = copy(dst, src)
	}
}

func BenchmarkCopyByAppend(b *testing.B) {
	src := make(Slice, 10240)
	src[1024] = 1
	dst := make(Slice, 10240)

	for i := 0; i < b.N; i++ {
		_ = CopyByAppend(src[:], src)
		_ = CopyByAppend(dst, src)
	}
}
