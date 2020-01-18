package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {}

// will stop cause v is copyed before range
func rangeFiniteLoop() {
	v := []int{1, 2, 3}
	for i := range v {
		v = append(v, i)
	}
}

func rangeModify() {
	arr := [2]int{1, 2}
	res := []*int{}
	for _, v := range arr {
		res = append(res, &v)
	}
	// expect: 1 2
	fmt.Println(*res[0], *res[1])
	// but output: 2 2

	// right way1:
	res = []*int{}
	for _, v := range arr {
		// use local variable
		v := v
		res = append(res, &v)
	}
	// right way2:
	res = []*int{}
	for k := range arr {
		res = append(res, &arr[k])
	}
}

// Use Array Pointers as Arrays
func rangeBigArr() {
	var arr [102400]int
	// not reconmend
	for i, n := range arr {
		_ = i
		_ = n
	}
	// reconmend &arr
	for i, n := range &arr {
		_ = i
		_ = n + 1
	}
	// reconmend arr[:]
	for i, n := range arr[:] {
		_ = i
		_ = n
	}
}

// nil array pointer cant iterate
func rangeNilPointerArray() {
	var p *[5]int // nil

	for i := range p { // okay
		fmt.Println(i)
	}

	for i := range p { // okay
		fmt.Println(i)
	}

	for i, n := range p { // panic
		fmt.Println(i, n)
	}
}

// go tool compile -S for-range.go |grep memclr
// also can check go src code: src/cmd/compile/internal/gc/range.go
// The memclr Optimization
func rangeResetOptimization() {
	a := []int{1, 2, 3, 4, 5}
	// array loop reset has optimize
	for i := range a {
		a[i] = 0
	}
	// slice loop reset has optimize
	s := a[:]
	for i := range s {
		s[i] = 0
	}

	// array pointer loop reset has [not] optimize
	arr := [10]int{1, 2, 4}
	for i := range &arr {
		a[i] = 0
	}
}

func rangeMapWhileCreateElem() {
	var createElemDuringIterMap = func() {
		var m = map[int]int{1: 1, 2: 2, 3: 3}
		for i := range m {
			m[4] = 4
			fmt.Printf("%d%d ", i, m[i])
		}
	}
	for i := 0; i < 50; i++ {
		//some line will not show 44, some line will
		createElemDuringIterMap()
		fmt.Println()
	}
}

func rangeMapWhileDeleteElem() {
	var m = map[int]int{1: 1, 2: 2, 3: 3}
	//only del key once, and not del the current iteration key
	var o sync.Once
	for i := range m {
		o.Do(func() {
			for _, key := range []int{1, 2, 3} {
				if key != i {
					fmt.Printf("when iteration key %d, del key %d\n", i, key)
					delete(m, key)
					break
				}
			}
		})
		fmt.Printf("%d%d ", i, m[i])
	}
}

func rangeGoroutineClosure() {
	// wrong
	var m = []int{1, 2, 3}
	for i := range m {
		go func() {
			fmt.Print(i)
		}()
	}

	// output: 222

	// right way1:
	for i := range m {
		// pass param
		go func(i int) {
			fmt.Print(i)
		}(i)
	}

	// right way2:
	for i := range m {
		// use local variable
		i := i
		go func() {
			fmt.Print(i)
		}()
	}
	// block 1ms to wait goroutine finished
	time.Sleep(time.Millisecond)
}
