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
Aligned    0.42ns ±22%
UnAligned  0.43ns ±19%
*/
type UnAligned struct {
	b   [7]byte
	i64 [8]byte
}
type Aligned struct {
	b   [8]byte
	i64 [8]byte
}

func BenchmarkAligned(b *testing.B) {
	x := Aligned{}
	tmp := (*int64)(unsafe.Pointer(&x.i64))
	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
func BenchmarkUnAligned(b *testing.B) {
	x := UnAligned{}
	tmp := (*int64)(unsafe.Pointer(&x.i64))
	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
