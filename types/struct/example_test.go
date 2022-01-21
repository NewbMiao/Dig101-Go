package main_test

import "fmt"

type T1 struct {
	a struct{}
	x int64
}

type T2 struct {
	x int64
	a struct{} // final zero field
}

type (
	T3 struct{ a struct{} }
	T4 struct{}
)

// For final zero field, see issue: https://github.com/golang/go/issues/9401
// 	and commit: https://go-review.googlesource.com/c/go/+/33475/2/src/reflect/type.go#2587
// nolint: govet // example
func ExampleStructEmpty() {
	var b1, b2 T4
	a1 := &T1{}
	a2 := &T2{}
	a3 := &T3{}
	fmt.Printf("a1: %v \ta2: %v \ta3: %v \tb1: %v \tb2: %v\n", a1, a2, a3, b1, b2)
	// &a3.a, &b1, &b2 point to zerobase, while &a1.a, &a2.a not
	fmt.Printf("&a1.a: %p \t&a2.a: %p \t&a3.a: %p \t&b1: %p \t&b2: %p\n", &a1.a, &a2.a, &a3.a, &b1, &b2)

	// Output like:
	// a1: &{{} 0} 			a2: &{0 {}} 			a3: &{{}} 			b1: {} 			b2: {}
	// &a1.a: 0xc0000162d0 	&a2.a: 0xc0000162e8 	&a3.a: 0x12501d8 	&b1: 0x12501d8 	&b2: 0x12501d8
}
