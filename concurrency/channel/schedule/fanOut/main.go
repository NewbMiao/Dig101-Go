package main

import (
	"fmt"
	"reflect"
	"time"
)

// nolint: gocognit // example
func fanOut(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() { // 退出时关闭所有的输出chan
			for i := range out {
				close(out[i])
			}
		}()

		for v := range ch { // 从输入chan中读取数据
			v := v
			for i := range out {
				i := i
				if !async {
					out[i] <- v // 放入到输出chan中，同步方式
					continue
				}
				go func() { // 异步,避免一个out阻塞的时候影响其他out
					out[i] <- v
				}()
			}
		}
	}()
}

func fanOutInReflect(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() { // 退出时关闭所有的输出chan
			for i := range out {
				close(out[i])
			}
		}()
		// 构造SelectCase slice
		cases := []reflect.SelectCase{
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			},
		}
		// 循环，从cases中选择一个可用的
		_, v, ok := reflect.Select(cases)
		if !ok { // 此channel已经close
			return
		}
		for i := range out {
			out[i] <- v.Interface()
		}
	}()
}

func main() {
	out := make([]chan interface{}, 3)
	for i := range out {
		out[i] = make(chan interface{})
	}
	ch := make(chan interface{})
	go func() {
		ch <- 1
	}()
	time.Sleep(time.Millisecond)
	fanOutInReflect(ch, out, true)
	for i := range out {
		fmt.Println("got: ", <-out[i])
	}
}
