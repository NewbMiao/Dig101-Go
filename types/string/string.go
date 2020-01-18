package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
	"unicode/utf8"
)

func main() {
	log.SetFlags(0)
	log.Printf("Type %#v:", reflect.StringHeader{})

	// compare is about underlying bytes sequence
	longStringCmp()

	// concat string avoid copy ( + is recomend way for concat, see concat_test's benchmark)
	concatStringWithCopyUnderlyingBytes()

	// 3way for-range
	forRangeString()
}

func longStringCmp() {
	bs := make([]byte, 1<<26) //64MB
	s0 := string(bs)
	s1 := string(bs)
	s2 := s1

	// s0, s1 and s2 are three equal strings.
	// The underlying bytes of s0 is a copy of bs.
	// The underlying bytes of s1 is also a copy of bs.
	// The underlying bytes of s0 and s1 are two
	// different copies of bs.
	// s2 shares the same underlying bytes with s1.

	startTime := time.Now()
	_ = s0 == s1
	duration := time.Since(startTime)
	log.Println("duration for (s0 == s1):", duration)

	startTime = time.Now()
	_ = s1 == s2
	duration = time.Since(startTime)
	log.Println("duration for (s1 == s2):", duration)
	log.Println("1ms is 1000000ns! So please try to avoid comparing two long strings if they don't share the same underlying byte sequence.")

}

var s string
var x = []byte{1023: 'x'}
var y = []byte{1023: 'y'}

func fc() {
	// None of the below 4 conversions will
	// copy the underlying bytes of x and y.
	// Surely, the underlying bytes of x and y will
	// be still copied in the string concatenation.
	if string(x) != string(y) {
		s = (" " + string(x) + string(y))[1:]
	}
}

func fd() {
	// Only the two conversions in the comparison
	// will not copy the underlying bytes of x and y.
	if string(x) != string(y) {
		// Please note the difference between the
		// following concatenation and the one in fc.
		s = string(x) + string(y)
	}
}

func concatStringWithCopyUnderlyingBytes() {
	/*
		https://go101.org/article/string.html
		Compiler Optimizations for Conversions Between Strings and Byte Slices: avoid the duplicate copies
		a conversion (from string to byte slice) which follows the range keyword in a for-range loop.
		a conversion (from byte slice to string) which is used as a map key in map element indexing syntax.
		a conversion (from byte slice to string) which is used in a comparison.
		a conversion (from byte slice to string) which is used in a string concatenation, and at least one of concatenated string values is a non-blank string constant.
	*/
	log.Println(`Way1: (" " + string(x) + string(y))[1:]`, testing.AllocsPerRun(1, fc)) // 1
	log.Println(`Way2: string(x) + string(y)`, testing.AllocsPerRun(1, fd))             // 3
}

func forRangeString() {
	s := "éक्षिaπ囧"
	log.Printf("string %s, len:%d len([]rune):%d, RuneCountInString:%d", s, len(s), len([]rune(s)), utf8.RuneCountInString(s))
	var sl string
	for i, rn := range s {
		sl += fmt.Sprintf("%2v: 0x%x %v ", i, rn, string(rn))
	}

	log.Println("range string: ", sl)

	sl = ""
	for i, bs := range []rune(s) {
		sl += fmt.Sprintf("%2v: 0x%x %v ", i, bs, string(bs))
	}
	log.Println("range runes:", sl)

	sl = ""
	for i, bs := range []byte(s) {
		sl += fmt.Sprintf("%2v: 0x%x %v ", i, bs, string(bs))
	}
	log.Println("range bytes:", sl)

	// above same as
	// for i := 0; i < len(s); i++ {
	// 	log.Printf("The byte at index %v: 0x%x %v\n", i, s[i], string(s[i]))
	// }

}
