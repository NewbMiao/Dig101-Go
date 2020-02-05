package main

func main() {

}

func unhashableType(){
	var m = map[interface{}]int{}
	var i interface{} = []int{}
	//panic: runtime error: hash of unhashable type []int
	println(m[i])
	//panic: runtime error: hash of unhashable type []int
	delete(m, i)
}

func unaddressable(){
	m0 := map[int]int{}
	// ❎ cannot take the address of m0[0]
	_ = &m0[0]

	m := make(map[int][2]int)
	// ✅
	m[0] = [2]int{1, 0}
	// ❎ cannot assign to m[0][0]
	m[0][0] = 1
	// ❎ cannot take the address of m[0]
	_ = &m[0]

	type T struct{ v int }
	ms := make(map[int]T)
	// ✅
	ms[0] = T{v: 1}
	// ❎ cannot assign to struct field ms[0].v in map
	ms[0].v = 1
	// ❎ cannot take the address of ms[0]
	_ = &ms[0]

}