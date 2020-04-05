package main

import (
	"sync/atomic"
	"testing"
	"unsafe"
)

/*
# This benchmark is useless now. Needs figure out a new way to do it.

# No gc,schedule,no optimize(escape and inline)
GOGC=off GODEBUG=asyncpreemptoff=1 go test -gcflags='-N -l' . -run none -bench . -count 20 -cpu 1 > b.txt && benchstat b.txt

name       time/op
UnAligned  3.58ns ± 9%
Aligned    3.56ns ±11%

# also can try use docker:
docker build -t  gobench-structalign .
docker run --rm   gobench-structalign
*/
const N = 256

type AItem [N]byte

func TestUnalignedAddress(t *testing.T) {
	x := AItem{}
	ptr := unsafe.Pointer(&x[9])
	intPtr := (*int64)(ptr)
	// equal to: unsafe.Pointer(uintptr(unsafe.Pointer(&x))+9)
	if uintptr(unsafe.Pointer(intPtr))%unsafe.Sizeof(intPtr) == 0 {
		t.Error("Not unaligned uintptr(ptr)")
	}
}
func TestAlignedAddress(t *testing.T) {
	x := AItem{}
	ptr := unsafe.Pointer(&x[8])
	intPtr := (*int64)(ptr)
	// equal to: unsafe.Pointer(uintptr(unsafe.Pointer(&x))+8)
	if uintptr(unsafe.Pointer(intPtr))%unsafe.Sizeof(intPtr) != 0 {
		t.Error("Not unaligned uintptr(ptr)")
	}
}

func BenchmarkUnAligned(b *testing.B) {
	for i := 1; i < b.N; i++ {
		accessAddr(1)
	}
}
func BenchmarkAligned(b *testing.B) {
	for i := 0; i < b.N; i++ {
		accessAddr(0)
	}
}

// access multiple time with aligned or unaligned address
func accessAddr(start int) {
	type VType [8]byte
	offset := int(unsafe.Alignof(VType{}))
	for i := start; i < N; i += offset {
		x := AItem{N - 1: 1}
		// avoid using uintptr -> unsafe.Pointer, the value that uintptr pointed to maybe lost.
		// tmp := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + 9))
		ptr := unsafe.Pointer(&x[i])
		tmp := VType{1: 1}
		atomic.StorePointer(&ptr, unsafe.Pointer(&tmp))
		_ = *((*VType)(ptr))
	}
}
