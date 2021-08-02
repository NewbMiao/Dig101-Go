package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var (
	sf    = singleflight.Group{}
	wg    sync.WaitGroup
	count = 20
)

func main() {
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			do, err, shared := sf.Do("number", Request)
			fmt.Println(i, " resp: ", do, " err: ", err, " shared: ", shared)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func Request() (interface{}, error) {
	return "process time: " + time.Now().String(), nil
}
