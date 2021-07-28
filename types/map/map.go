package main

import "math"

func main() {
}

// slice is incomparable and cant use for hash
// incomparable type:
// func, map, slice
// any array or struct contains func/map/slice.
func unhashableType() {
	m := map[interface{}]int{}
	var i interface{} = []int{}
	// panic: runtime error: hash of unhashable type []int
	println(m[i])
	// panic: runtime error: hash of unhashable type []int
	delete(m, i)
}

// can not take address of map elem.
func unaddressable() {
	// m0 := map[int]int{}
	// ❎ cannot take the address of m0[0]
	// _ = &m0[0]

	m := make(map[int][2]int)
	// ✅
	m[0] = [2]int{1, 0}
	// ❎ cannot assign to m[0][0]
	// m[0][0] = 1
	// ❎ cannot take the address of m[0]
	// _ = &m[0]

	type T struct{ v int }
	ms := make(map[int]T)
	// ✅
	ms[0] = T{v: 1}
	// ❎ cannot assign to struct field ms[0].v in map
	// ms[0].v = 1
	// ❎ cannot take the address of ms[0]
	// _ = &ms[0]
}

// NaN: key != key, use NaN as key make no sense.
func unreachable() {
	n1, n2 := math.NaN(), math.NaN()
	m := map[float64]int{}
	m[n1], m[n2] = 1, 2
	println(n1 == n2, m[n1], m[n2])
	// output: false 0 0
}
