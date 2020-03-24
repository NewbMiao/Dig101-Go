package main

// struct for demo

// structlayout -json github.com/NewbMiao/Dig101-Go Ag|structlayout-svg -t "align-guarantee" > ag.svg
type Ag struct {
	arr [2]int64 // 16
	sl  []int64  // 24
	bl  bool     // 1 padding 7
	ptr *int64   // 8
	st  struct { // 16
		str string
	}
}

// structlayout -json github.com/NewbMiao/Dig101-Go tooMuchPadding|structlayout-optimize -r
type tooMuchPadding struct {
	i16 int16
	i64 int64
	i8  int8
	i32 int32
	ptr *string
	b   bool
}

// tooMuchPadding optimized
type optimized struct {
	i64 int64
	ptr *string
	i32 int32
	i16 int16
	i8  int8
	b   bool
}

// 64word align gurrantee on 32-bit arch
// GOARCH=386 structlayout -json github.com/NewbMiao/Dig101-Go c2 | structlayout-svg -t "int64 first field" > i64_first.svg
type c2 struct {
	val   int64 // pos 0
	val2  int64 // pos 8
	valid bool  // pos 16
}

type T struct {
	val2 int64
	_    int16
}
type c3 struct {
	val   T
	valid bool
}

type c4 struct {
	val   int64 // pos 0
	valid bool  // pos 8
	// 或者 _ uint32
	_    [4]byte // pos 9; to correct padding one more 4bytes
	val2 int64   // pos 16
}

type c5 struct {
	val   int64
	valid bool
	// the first element in slices of 64-bit
	// elements will be correctly aligned
	// 此处切片相当指针，数据是指向底层开辟的64位字数组，如c1
	val2 []int64
}

type c51 struct {
	val   int64
	valid bool
	// the first element in slices of 64-bit
	// elements will be correctly aligned
	// 此处切片相当指针，数据是指向底层开辟的64位字数组，如c1
	val2 [2]int64
}

type c6 struct {
	val   int64
	valid bool
	val2  *int64
}

func main()  {
	
}
