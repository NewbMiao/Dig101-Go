package main

import (
	"testing"
	"unsafe"
)

/*
go get golang.org/x/perf/cmd/benchstat

GOARCH=386 go test github.com/NewbMiao/Dig101-Go/types/struct -bench . -count 10 > old.txt

// 386_amd64
benchstat old.txt
name          time/op
UnAligned-12  0.82ns ± 6%
Aligned-12    0.52ns ± 1%
*/
var ptrSize uintptr

func init() {
	ptrSize = unsafe.Sizeof(uintptr(1))
}

func BenchmarkUnAligned(b *testing.B) {
	type UnAligned struct {
		b [25]byte
	}
	x := UnAligned{}
	address := uintptr(unsafe.Pointer(&x.b))
	if address%ptrSize == 0 {
		b.Error("Not unaligned address")
	}
	tmp := (*int64)(unsafe.Pointer(&x.b))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
func BenchmarkAligned(b *testing.B) {
	type Aligned struct {
		b [24]byte
	}

	x := Aligned{}
	address := uintptr(unsafe.Pointer(&x.b))
	if address%ptrSize != 0 {
		b.Error("Not aligned address")
	}
	tmp := (*int64)(unsafe.Pointer(&x.b))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		*tmp = int64(100)
	}
}
