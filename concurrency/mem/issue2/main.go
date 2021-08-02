package main

var (
	a    string
	done bool
)

func setup() {
	a = "hello, world"
	done = true // 指令重排后，done 和 a 赋值顺序不能保证
}

func main() {
	go setup()
	for !done { // 可能会一直hang在这
	}
	print(a)
}
