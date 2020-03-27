package main

import (
	"testing"
	"unsafe"
)

/*
perflock go test -gcflags='-N -l' github.com/NewbMiao/Dig101-Go/types/struct -bench . -count 3 > old.txt
benchstat old.txt

name         time/op
UnAligned-6  1.87ns ± 5%
Aligned-6    1.47ns ± 2%

// also can try use docker:
docker build -t  gobench-structalign https://raw.githubusercontent.com/NewbMiao/Dig101-Go/master/types/struct/Dockerfile
docker run --rm   gobench-structalign
*/
var ptrSize uintptr

func init() {
	ptrSize = unsafe.Sizeof(uintptr(1))
}

type SType struct {
	b [32]byte
}

func BenchmarkUnAligned(b *testing.B) {
	x := SType{}
	address := uintptr(unsafe.Pointer(&x.b)) + 1
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
	x := SType{}
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
