package main

import (
	"fmt"
	"reflect"
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
		default: // 超过两个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-OrRecur(channels[:m]...):
			case <-OrRecur(channels[m:]...):
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

func OrInReflect(channels ...<-chan interface{}) <-chan interface{} {
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

	<-OrRecur(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)
	fmt.Printf("[orDone] done after %v\n", time.Since(start))

	start = time.Now()

	<-OrInReflect(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)
	fmt.Printf("[orDone in reflect] done after %v", time.Since(start))
}
