package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

func fanOut(ch <-chan interface{}, out []chan interface{}) {
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
				out[i] <- v // 放入到输出chan中，同步方式
			}
		}
	}()
}

func fanOutAsync(ch <-chan interface{}, out []chan interface{}) {
	go func() {
		var wg sync.WaitGroup
		defer func() { // 退出时关闭所有的输出chan
			wg.Wait()
			for i := range out {
				close(out[i])
			}
		}()

		for v := range ch { // 从输入chan中读取数据
			v := v
			for i := range out {
				i := i
				wg.Add(1)
				go func() { // 异步,避免一个out阻塞的时候影响其他out
					out[i] <- v
					wg.Done()
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
	fmt.Println("任务进行中，当前协程数:", runtime.NumGoroutine())
	out := make([]chan interface{}, 3)
	for i := range out {
		out[i] = make(chan interface{})
	}
	done := make(chan struct{})
	ch := generator.AsStream(done, []interface{}{1, 2, 3})
	fanOutAsync(ch, out)
	for i := range out {
		fmt.Println("got: ", <-out[i])
	}
	fmt.Println("任务结束，当前协程数:", runtime.NumGoroutine())
}
