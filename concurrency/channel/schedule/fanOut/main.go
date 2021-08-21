package main

import (
	"fmt"
	"reflect"
	"time"
)

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

func fanOutReflect(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		defer func() { // 退出时关闭所有的输出chan
			for i := range out {
				close(out[i])
			}
		}()
		cases := make([]reflect.SelectCase, len(out))
		// 构造SelectCase slice
		for i := range cases {
			cases[i].Dir = reflect.SelectSend
		}
		for v := range ch {
			v := v
			for i := range cases {
				cases[i].Chan = reflect.ValueOf(out[i])
				cases[i].Send = reflect.ValueOf(v)
			}
			for range cases {
				chosen, _, _ := reflect.Select(cases)
				// 已发送过，用nil阻塞，避免再次发送
				cases[chosen].Chan = reflect.ValueOf(nil)
			}
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
	fanOutReflect(ch, out)
	for i := range out {
		fmt.Println("got: ", <-out[i])
	}
}
