package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var shortStr = string(make([]byte, 16))
var longStr = string(make([]byte, 1024))

// go test . -bench  BenchmarkConcat -benchtime=3s -benchmem

// The standard Go compiler makes optimizations for string concatenations by using the + operator.
// So generally, using + operator to concatenate strings is convenient and efficient
// if the number of the concatenated strings is known at compile time.
// otherwise byte.Buffer will be more efficient
func BenchmarkConcatLess(b *testing.B) {
	concatIter(b, shortStr, 1)
	concatIter(b, longStr, 1)
}

func BenchmarkConcatMore(b *testing.B) {
	concatIter(b, shortStr, 100)
	concatIter(b, longStr, 100)
}

func concatIter(b *testing.B, v string, iter int) {
	var k string
	if len(v) > 32 {
		k += "long"
	} else {
		k += "short"
	}

	// implements a Go string concatenation x+y+z+...
	// is small than 32 would use [32]byte to concat
	// if not alloc memory to concat
	b.Run(k+"Plus", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var s string
			for i := 0; i < iter; i++ {
				s += "my_string" + v
				b.SetBytes(int64(len(v)))
			}

		}
	})

	// initial a []byte buffer which lenth is 64, and grow double when need,
	// otherwise just reslice itself
	b.Run(k+"Buffer", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var bf bytes.Buffer
			for i := 0; i < iter; i++ {
				bf.WriteString("my_string")
				bf.WriteString(v)
				b.SetBytes(int64(len(v)))
			}
		}
	})

	// inside use strings.Builder
	b.Run(k+"Join", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var s string
			for i := 0; i < iter; i++ {
				s += strings.Join([]string{"my_string", v}, "")
				b.SetBytes(int64(len(v)))
			}
		}
	})

	// minimizes memory copying, contains copyCheck
	b.Run(k+"Builder", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var bf strings.Builder
			for i := 0; i < iter; i++ {
				bf.WriteString("my_string")
				bf.WriteString(v)
				b.SetBytes(int64(len(v)))
			}
		}
	})

	// use simple []byte
	b.Run(k+"Sprint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var s string
			for i := 0; i < iter; i++ {

				s += fmt.Sprint("my_string", v)
				b.SetBytes(int64(len(v)))
			}
		}
	})

	// use simple []byte and need do formatLoop
	b.Run(k+"Sprintf", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var s string
			for i := 0; i < iter; i++ {

				s += fmt.Sprintf("my_string%s", v)
				b.SetBytes(int64(len(v)))
			}
		}
	})
}
