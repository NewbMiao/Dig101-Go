package main

import (
	"testing"
	"unsafe"
)

/*
go get golang.org/x/perf/cmd/benchstat

go test github.com/NewbMiao/Dig101-Go/types/struct -bench .  -count=10 -cpu 1 > old.txt

benchstat old.txt                                                                      
name       time/op
Aligned    0.43ns ±32%
UnAligned  0.44ns ±34%
*/
type UnAligned struct {
	b [15]byte
}
type Aligned struct {
	b [16]byte
}

func BenchmarkAligned(b *testing.B) {
	x := Aligned{}
	// fmt.Println(uintptr(unsafe.Pointer(&x)) % 8) // == 0

	tmp := (*int64)(unsafe.Pointer(&x))
	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
func BenchmarkUnAligned(b *testing.B) {
	x := UnAligned{}
	// fmt.Println(uintptr(unsafe.Pointer(&x)) % 8) // == 1

	tmp := (*int64)(unsafe.Pointer(&x))
	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
