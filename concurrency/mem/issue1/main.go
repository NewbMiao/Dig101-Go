package main

import (
	"fmt"
	"time"
)

var a, b int

func f() {
	a = 1 // w之前的写操作
	b = 2 // 写操作w
	fmt.Printf("g1: inside f, ba: %d %d\n", b, a)
}

func g() {
	fmt.Printf("g2: inside g, ba: %d %d\n", b, a)

	// print(b) // 读操作r
	// print(a) // ???
}

func main() {
	go f() // g1
	g()    // g2
	g()
	time.Sleep(time.Microsecond)
}
