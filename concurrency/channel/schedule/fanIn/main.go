package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/schedule/generator"
)

func main() {
	fmt.Println("任务进行中，当前协程数:", runtime.NumGoroutine())
	done := make(chan struct{})
	defer close(done)
	streams := make([]<-chan interface{}, 100)
	for i := range streams {
		streams[i] = generator.AsStream(done, []int{1, 2, 3})
	}
	ch := fanIn(streams...)
	for range ch {
	}
	fmt.Println("任务结束，当前协程数:", runtime.NumGoroutine())
}

func fanIn(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))
		for _, ch := range chans {
			go func(ch <-chan interface{}) {
				for v := range ch {
					out <- v
				}
				wg.Done()
			}(ch)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func fanInReflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选择一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok { // 此channel已经close
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

func fanInRecur(chans ...<-chan interface{}) <-chan interface{} {
	switch len(chans) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(
			fanInRecur(chans[:m]...),
			fanInRecur(chans[m:]...))
	}
}

func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil { // 只要还有可读的chan
			select {
			case v, ok := <-a:
				if !ok { // a 已关闭，设置为nil
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok { // b 已关闭，设置为nil
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}
