package water

import (
	"context"

	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
)

// 定义水分子合成的辅助数据结构.
type H2O struct {
	semaH *semaphore.Weighted         // 氢原子的信号量
	semaO *semaphore.Weighted         // 氧原子的信号量
	b     cyclicbarrier.CyclicBarrier // 循环栅栏，用来控制合成
}

func New() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(2), // 氢原子需要两个
		semaO: semaphore.NewWeighted(1), // 氧原子需要一个
		b:     cyclicbarrier.New(3),     // 需要三个原子才能合成
	}
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)

	releaseHydrogen()                 // 输出H
	h2o.b.Await(context.Background()) // 等待栅栏放行
	h2o.semaH.Release(1)              // 释放氢原子空槽
}

func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)

	releaseOxygen()                   // 输出H
	h2o.b.Await(context.Background()) // 等待栅栏放行
	h2o.semaO.Release(1)              // 释放氢原子空槽
}
