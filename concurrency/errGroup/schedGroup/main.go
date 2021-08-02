package main

import (
	"context"
	"log"
	"time"

	"github.com/mdlayher/schedgroup"
)

func main() {
	sg := schedgroup.New(context.Background())
	// 设置子任务分别在100、200、300之后执行
	for i := 0; i < 3; i++ {
		n := i + 1
		sg.Delay(time.Duration(n*100)*time.Millisecond, func() {
			log.Println(n) // 输出任务编号
		})
	}

	// 等待所有的子任务都完成
	if err := sg.Wait(); err != nil {
		log.Fatalf("failed to wait: %v", err)
	}
}
