package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var shortStr = "my_string" + string(make([]byte, 16))
var longStr = "my_string" + string(make([]byte, 256))

// go test . -bench  BenchmarkConcat -benchtime=3s -benchmem

func BenchmarkConcatShort10(b *testing.B) {
	concatIter(b, shortStr, 10)
}

func BenchmarkConcatShort1000(b *testing.B) {
	concatIter(b, shortStr, 1000)
}

func BenchmarkConcatLong10(b *testing.B) {
	concatIter(b, longStr, 10)
}

func BenchmarkConcatLong1000(b *testing.B) {
	// concatIter(b, longStr, 1000)
	concatIterWithGrowInit(b, longStr, 1000)
}

// grow init avoid realloc cost
func concatIterWithGrowInit(b *testing.B, v string, iter int) {
	b.Run("BufferInit", func(b *testing.B) {
		var bf bytes.Buffer
		bf.Grow(iter * len(v))
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			bf.Reset()
			for i := 0; i < iter; i++ {
				bf.WriteString(v)
			}
			_ = bf.String()
		}

	})

	b.Run("BuilderInit", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var bf strings.Builder
			bf.Grow(iter * len(v))
			for i := 0; i < iter; i++ {
				bf.WriteString(v)
			}
			_ = bf.String()
		}

	})
}

func concatIter(b *testing.B, v string, iter int) {
	// implements a Go string concatenation x+y+z+...
	// is small than 32 would use [32]byte to concat
	// if not alloc memory to concat
	b.Run("Plus", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var s string
			for i := 0; i < iter; i++ {
				s += v
			}
			_ = s
		}

	})

	// initial a []byte buffer which lenth is 64, and grow double when need,
	// otherwise just reslice itself
	b.Run("Buffer", func(b *testing.B) {
		var bf bytes.Buffer
		for n := 0; n < b.N; n++ {
			bf.Reset()
			for i := 0; i < iter; i++ {
				bf.WriteString(v)
			}
			_ = bf.String()
		}

	})

	// inside use strings.Builder
	b.Run("Join", func(b *testing.B) {

		strsArgs := make([]string, iter)
		for i := 0; i < iter; i++ {
			strsArgs[i] = v
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			_ = strings.Join(strsArgs, "")
		}
	})

	// minimizes memory copying, contains copyCheck, and when want string output will more efficient than buffer
	b.Run("Builder", func(b *testing.B) {
		var bf strings.Builder
		for n := 0; n < b.N; n++ {
			bf.Reset()
			for i := 0; i < iter; i++ {
				bf.WriteString(v)
			}
			_ = bf.String()
		}

	})

	// use simple []byte
	b.Run("Sprint", func(b *testing.B) {
		sprintArgs := make([]interface{}, iter)
		for i := 0; i < iter; i++ {
			sprintArgs[i] = v
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = fmt.Sprint(sprintArgs...)
		}
	})

	// use simple []byte and need do formatLoop
	b.Run("Sprintf", func(b *testing.B) {
		sprintArgs := make([]interface{}, iter)
		for i := 0; i < iter; i++ {
			sprintArgs[i] = v
		}
		sprintFmt := strings.Repeat("%s", iter)

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			_ = fmt.Sprintf(sprintFmt, sprintArgs...)
		}
	})
}
