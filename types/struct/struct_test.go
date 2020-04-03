package main

import (
	"testing"
	"unsafe"
)

/*
# This benchmark is useless now. Needs figure out a new way to do it.

# No gc,schedule,no optimize(escape and inline)
GOGC=off GODEBUG=asyncpreemptoff=1 go test -gcflags='-N -l' . -bench . -count 20 > old.txt
benchstat old.txt

name         time/op
UnAligned-6  1.82ns ± 0%
Aligned-6    1.82ns ± 0%

# also can try use docker:
docker build -t  gobench-structalign .
docker run --rm   gobench-structalign
*/
var ptrSize uintptr

func init() {
	ptrSize = unsafe.Sizeof(uintptr(1))
}

type SType struct {
	b [64]byte
}

func BenchmarkUnAligned(b *testing.B) {
	x := SType{}

	ptr := unsafe.Pointer(&x.b[9])
	// equal to: unsafe.Pointer(uintptr(unsafe.Pointer(&x.b))+9)
	if uintptr(ptr)%ptrSize == 0 {
		b.Error("Not unaligned uintptr(ptr)")
	}
	// avoid using uintptr -> unsafe.Pointer, the value that uintptr pointed to maybe lost.
	// tmp := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&x.b)) + 9))
	tmp := (*int64)(ptr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
func BenchmarkAligned(b *testing.B) {
	x := SType{}
	ptr := unsafe.Pointer(&x.b[8])
	// equal to: unsafe.Pointer(uintptr(unsafe.Pointer(&x.b))+8)
	if uintptr(ptr)%ptrSize != 0 {
		b.Error("Not aligned uintptr(ptr)")
	}
	tmp := (*int64)(ptr)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		*tmp = int64(i)
	}
}
