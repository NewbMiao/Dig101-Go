package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"unsafe"
)

func main() {
	alignedAddress()
	shadowField()
	safeAtomicAccess64bitWordOn32bitArch()
	alignPadding()
	zeroField()
}

/*
M
https://github.com/golang/go/issues/19149

Output will like this, show address is 4bytes aligned, but not consist after mod 8
$ GOARCH=386 go run main.go
290566072 0 290566060 4 290566048 0

$ GOARCH=amd64 go run main.go
824635031404 4 824635031392 0 824635031380 4.

*/
type M struct {
	x [3]uint32
}

func alignedAddress() {
	var a, b, c M

	println(
		uintptr(unsafe.Pointer(&a.x)), uintptr(unsafe.Pointer(&a.x))%8,
		uintptr(unsafe.Pointer(&b.x)), uintptr(unsafe.Pointer(&b.x))%8,
		uintptr(unsafe.Pointer(&c.x)), uintptr(unsafe.Pointer(&c.x))%8,
	)
}

// use max align of struct fields
// 64-bit arch: 8-byte aligned; 32-bit arch: 4-byte aligned
// if has gap between fields use padding.
func alignPadding() {
	type T1 struct {
		a [2]int8
		b int64
		c int16
	}

	type T2 struct {
		a [2]int8
		c int16
		b int64
	}

	fmt.Printf("arrange fields to reduce size:\n"+
		"T1 align: %d, size: %d\n"+
		"T2 align: %d, size: %d\n",
		unsafe.Alignof(T1{}), unsafe.Sizeof(T1{}),
		unsafe.Alignof(T2{}), unsafe.Sizeof(T2{}))
}

func zeroField() {
	type T1 struct {
		a struct{}
		x int64
	}

	type T2 struct {
		x int64
		// pad bytes avoid memory leak when use address of this final zero field
		a struct{}
	}

	a1 := T1{}
	a2 := T2{}
	fmt.Printf("zero size struct{} in field:\n"+
		"T1 (not as final field) size: %d\n"+
		"T2 (as final field) size: %d\n",
		unsafe.Sizeof(a1), unsafe.Sizeof(a2)) // 16
}

/*
https://golang.org/pkg/sync/atomic/#pkg-note-BUG

On x86-32, the 64-bit functions use instructions unavailable before the Pentium MMX.

On non-Linux ARM, the 64-bit functions use instructions unavailable before the ARMv6k core.

On ARM, x86-32, and 32-bit MIPS, it is the caller's responsibility to arrange for 64-bit
alignment of 64-bit words accessed atomically.
The first word in a variable or in an allocated struct, array, or slice can be relied upon
to be 64-bit aligned.

https://go101.org/article/memory-layout.html#size-and-padding
https://stackoverflow.com/a/51012703/4431337

GOARCH=386 go run types/struct/struct.go.
*/
func safeAtomicAccess64bitWordOn32bitArch() {
	fmt.Println("32位系统下可原子安全访问的64位字：")

	var c0 int64

	fmt.Println("64位字本身：",
		atomic.AddInt64(&c0, 1))

	c1 := [5]int64{}
	fmt.Println("64位字数组、切片:",
		atomic.AddInt64(&c1[:][0], 1))

	c2 := struct {
		val   int64 // pos 0
		val2  int64 // pos 8
		valid bool  // pos 16
	}{}
	fmt.Println("结构体首字段为对齐的64位字及相邻的64位字:",
		atomic.AddInt64(&c2.val, 1),
		atomic.AddInt64(&c2.val2, 1))

	type T struct {
		val2 int64
		_    int16
	}

	c3 := struct {
		val   T
		valid bool
	}{}
	fmt.Println("结构体中首字段为嵌套结构体，且其首元素为64位字:",
		atomic.AddInt64(&c3.val.val2, 1))

	c4 := struct {
		val   int64 // pos 0
		valid bool  // pos 8
		// 或者 _ uint32
		_    [4]byte // pos 9; to correct padding one more 4bytes
		val2 int64   // pos 16
	}{}
	fmt.Println("结构体增加填充使对齐的64位字:",
		atomic.AddInt64(&c4.val2, 1))

	c5 := struct {
		val   int64
		valid bool
		// the first element in slices of 64-bit
		// elements will be correctly aligned
		// 此处切片相当指针，数据是指向底层开辟的64位字数组，如c1
		val2 []int64
	}{val2: []int64{0}}
	fmt.Println("结构体中64位字切片:",
		atomic.AddInt64(&c5.val2[0], 1))

	// 如果换成数组则会panic，
	// 因为结构体的数组的对齐还是依赖于结构体内字段
	// c51 := struct {
	//	val   int64
	//	valid bool
	//	val2  [3]int64
	// }{val2: [3]int64{0}}
	// fmt.Println("结构体中64位字切片:",
	//	atomic.AddInt64(&c51.val2[0], 1))

	c6 := struct {
		val   int64
		valid bool
		val2  *int64
	}{val2: new(int64)}
	fmt.Println("结构体中64位字指针:",
		atomic.AddInt64(c6.val2, 1))
}

func shadowField() {
	type Embedded struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	type Top struct {
		A interface{} `json:"a"`
		// 初始化为空结构体
		Embedded
		// 类似： *Embedded 初始化为空指针 nil
	}

	var a Top
	// B字段被提升，可直接操作
	a.B = "0"
	a.A = 1
	// 不同于a.Embedded.A

	jsonBytes := []byte(`{"a":["1","2"],"b":"3","a1":"4"}`)

	// 所有 json tag 会被提升用来匹配，遵循同名覆盖规则
	// 所以 tag "a" 会解码到 a.A
	json.Unmarshal(jsonBytes, &a)
	fmt.Printf("unmarshal json with embedded field: %+v\n", a)
	// {A:[1 2] Embedded:{A: B:3}}
	// 若Embedded.A json tag 为a1， 则输出为
	// {A:[1 2] Embedded:{A:4 B:3}}
}
