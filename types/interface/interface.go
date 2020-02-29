package main

import (
	"strconv"
)

type Stringer interface {
	String() string
}

// declare
func efaceAndiface() {
	var a interface{}
	var b Stringer
	println(a, b)
}

type Binary uint64

func (i Binary) String() string {
	return strconv.Itoa(int(i))
}

// conversion
func conversion() {
	var b Stringer
	var i Binary = 1
	b = i //convT64

	_ = b.String()
}

// GOSSAFUNC=main go1.14 build types/interface/interface.go
// to check virtual version of ssa: types/interface/ssa.html
func devirt() {
	var b Stringer = Binary(1)
	_ = b.String() //static call Binary.String
}

// typeAssert
func typeAssert() {
	var b interface{} = Binary(1)
	v, ok := b.(Stringer) //getitab
	println(v, ok)
}

// https://github.com/golang/go/wiki/InterfaceSlice
//func interfaceSlice() {
//	var dataSlice []int
//	var interfaceSlice []interface{} = dataSlice
//	_ = indirectiface{}
//}

func main() {
	efaceAndiface()
	conversion()
	devirt()
	typeAssert()
}
