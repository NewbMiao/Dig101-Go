package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.RWMutex
	// writer,稍微等待，然后制造一个调用Lock的场景
	go func() {
		time.Sleep(200 * time.Millisecond)
		mu.Lock()
		fmt.Print("Write Lock\n")
		time.Sleep(100 * time.Millisecond)
		mu.Unlock()
		fmt.Printf("Write Unlock\n")
	}()

	go func() {
		factorial(&mu, 10) // 计算10的阶乘, 10!
	}()

	select {}
}

// 递归调用计算阶乘. 环形依赖死锁.
func factorial(m *sync.RWMutex, n int) int {
	if n < 1 { // 阶乘退出条件
		return 0
	}

	fmt.Printf("RLock %d\n", n)
	m.RLock()
	defer func() {
		fmt.Printf("RUnlock %d\n", n)
		m.RUnlock()
	}()
	time.Sleep(100 * time.Millisecond)

	return factorial(m, n-1) * n // 递归调用
}
