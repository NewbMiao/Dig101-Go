package main

type T struct {
	msg string
}

var g *T

func setup() {
	t := new(T)
	t.msg = "hello, world" // 指令重排后，可能g有值，而msg未赋值
	g = t
}

func main() {
	go setup()
	for g == nil {
	}
	print(g.msg)
}
