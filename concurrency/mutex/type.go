package mutex

import "sync"

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
)

// 扩展一个Mutex结构.
type Mutex struct {
	sync.Mutex
}
