package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"sort"
)

type T interface{}

func init() {
	log.SetFlags(0)
}

func discardLog() {
	log.SetOutput(ioutil.Discard)
}
func genSliceA() []T {
	return []T{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
}
func main() {
	log.Printf("Type %#v:", reflect.SliceHeader{})

	//append --------------------------------
	appendSliceDiff()

	var a []T
	//copy --------------------------------
	a = nil
	copySlice(a, "way 1 is not match")
	a = make([]T, 0)
	copySlice(a, "way 2 is not match")
	a = append(a, "test")
	copySlice(a, "all way match")

	log.Println("\nUse slice:", genSliceA())
	//cut --------------------------------
	// to see difference of mem profile:
	// go test . -bench  BenchmarkCutSlice -benchmem   -memprofile=mem.out
	// go tool pprof -http=localhost:8080 mem.out
	cutSlice(genSliceA(), 1, 3)
	cutSlicePointer(genSliceA(), 1, 3)

	//del --------------------------------
	delSliceItem(genSliceA(), 3)
	delSlicePointerItem(genSliceA(), 3)

	//insert --------------------------------
	insertSliceItem(genSliceA(), 3, "33")
	insertSliceItemWithoutNewSlice(genSliceA(), 4, "44")

	//pop --------------------------------
	var pop, popShift T
	a = genSliceA()
	pop, a = a[len(a)-1], a[:len(a)-1]
	log.Printf("pop slice, val: %v, after: %v", pop, a)

	//pop shift--------------------------------
	a = genSliceA()
	popShift, a = a[0], a[1:]

	log.Printf("popShift slice, val: %v, after: %v", popShift, a)

	//reverse--------------------------------
	reverseSlice(genSliceA())

	//shuffle--------------------------------
	shuffleSlice(genSliceA())

	//batch--------------------------------
	batchSlice(genSliceA(), 3)

	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1}
	deduplicateSlice(in)
}

func copySlice(a []T, notes string) {
	log.Printf("Copy slice, Got a:(%p,%+v), notes:%v", a, a, notes)

	var b []T
	// way 1
	b = make([]T, len(a))
	copy(b, a)

	log.Printf("way 1: b(%p,%+v)", b, b)

	// way 2
	b = append([]T(nil), a...)
	log.Printf("way 2: b(%p,%+v)", b, b)

	// way 3
	b = append(a[:0:0], a...) // See https://github.com/go101/go101/wiki
	log.Printf("way 3: b(%p,%+v)", b, b)
}

func cutSlice(a []T, from, to int) []T {
	l := len(a)
	if from < 0 && from >= l && to > from && to >= l {
		panic("runtime error: slice bounds out of range")
	}
	tmp := append(a[:from], a[to:]...)
	log.Printf("cut slice to %v which index is %d~%d (after origin slice would be %v)", tmp, from, to-1, a)
	return tmp
}
func cutSlicePointer(a []T, from, to int) []T {
	l := len(a)
	if from < 0 && from >= l && to > from && to >= l {
		panic("runtime error: slice bounds out of range")
	}
	copy(a[from:], a[to:]) //overlap the item which index is from~to
	for i := l - to + from; i < l; i++ {
		a[i] = nil //mark pointer unuse
	}
	a = a[:l-to+from]
	log.Printf("cut slice  to %v which index is %d~%d ", a, from, to-1)
	return a
}

func delSliceItem(a []T, i int) []T {
	l := len(a)
	if i > l {
		panic("runtime error: slice bounds out of range")
	}
	// tmp := append(a[:i], a[i+1:]...)
	tmp := a[:i+copy(a[i:], a[i+1:])]
	log.Printf("del slice to %v which index is %d (after origin slice would be %v)", tmp, i, a)
	return tmp
}
func delSlicePointerItem(a []T, i int) []T {
	l := len(a)
	if i > l {
		panic("runtime error: slice bounds out of range")
	}
	copy(a[i:], a[i+1:])
	a[l-1] = nil
	tmp := a[:l-1]
	// Think: why not use below way
	// a[i] = nil
	// tmp := append(a[:i], a[i+1:]...)

	log.Printf("del slice to %v which index is %d (after origin slice would be: %v)", tmp, i, a)
	return tmp
}

func insertSliceItem(a []T, i int, val T) []T {
	l := len(a)
	if i > l {
		panic("runtime error: slice bounds out of range")
	}
	tmp := append(a[:i], append([]T{val}, a[i:]...)...)
	log.Printf("insert slice to %v which index is %d val is %v (after origin slice would be %v)", tmp, i, val, a)
	return tmp
}

func insertSliceItemWithoutNewSlice(a []T, i int, val T) []T {
	l := len(a)
	if i > l {
		panic("runtime error: slice bounds out of range")
	}
	a = append(a, 0 /* use the zero value of the element type */)
	copy(a[i+1:], a[i:])
	a[i] = val
	log.Printf("insert slice to %v which index is %d val is %v", a, i, val)
	return a
}

func reverseSlice(a []T) []T {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	//or
	// for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
	// 	a[left], a[right] = a[right], a[left]
	// }
	log.Printf("reverse slice to %v", a)

	return a
}

func shuffleSlice(a []T) []T {
	// same as math/rand.Shuffle
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	log.Printf("shuffle slice to %v", a)

	return a

}

func batchSlice(actions []T, batchSize int) (batches [][]T) {
	for batchSize < len(actions) {
		actions, batches = actions[batchSize:], append(batches, actions[0:batchSize:batchSize])
	}
	batches = append(batches, actions)
	log.Printf("batch slice to each group size is %d's batches %v ", batchSize, batches)
	return
}

func deduplicateSlice(in []int) {
sort.Ints(in)
sl := fmt.Sprintf("deduplicate slice %v (sorted)", in)

j := 0
for i := 1; i < len(in); i++ {
	if in[j] == in[i] {
		continue
	}
	j++
	in[j] = in[i]
}
result := in[:j+1]
	log.Printf("%s to %v (after origin slice(sorted) would be %v)", sl, result, in)
}
