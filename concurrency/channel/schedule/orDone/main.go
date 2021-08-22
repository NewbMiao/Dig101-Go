package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

func Or(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var once sync.Once
		for _, c := range chans {
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() { close(out) })
				case <-out:
				}
			}(c)
		}
	}()
	return out
}
func OrWithIssue(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: // 2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //超过两个，二分法递归处理
			/*
				3个时有无限递归的问题:
				    f(3)
				 f(2)  f(3)
					f(2)  f(3)
					   f(2)  f(3)
								...
			*/
			m := len(channels) / 2
			select {
			case <-OrWithIssue(append(channels[:m:m], orDone)...):
			case <-OrWithIssue(append(channels[m:], orDone)...):
			}
		}
	}()

	return orDone
}

func OrRecurSimple(channels ...<-chan interface{}) <-chan interface{} {

	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-OrRecurSimple(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}
func OrRecur(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: // 2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		case 3: // 3个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			}
		default: // 超过3个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-OrRecur(append(channels[:m:m], orDone)...):
			case <-OrRecur(append(channels[m:], orDone)...):
			}
		}
	}()

	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func OrReflection(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有0个或者1个
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		// 利用反射构建SelectCase
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 随机选择一个可用的case
		reflect.Select(cases)
	}()

	return orDone
}

func main() {
	start := time.Now()
	go func() {
		time.Sleep(time.Second / 2)
		fmt.Println("任务进行中，当前协程数:", runtime.NumGoroutine())
	}()
	<-OrWithIssue(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)
	fmt.Printf("[orDone] done after %v\n", time.Since(start))
	fmt.Println("任务结束，当前协程数:", runtime.NumGoroutine())
}
