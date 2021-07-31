package main

import (
	"fmt"
	"time"

	"github.com/NewbMiao/Dig101-Go/concurrency/channel/jobs"
)

func orderPrint1Simple() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)
	go func() {
		for {
			fmt.Println("I'm goroutine 1")
			time.Sleep(1 * time.Second)
			ch2 <- 1 // I'm done, you turn
			<-ch1
		}
	}()

	go func() {
		for {
			<-ch2
			fmt.Println("I'm goroutine 2")
			time.Sleep(1 * time.Second)
			ch3 <- 1
		}
	}()

	go func() {
		for {
			<-ch3
			fmt.Println("I'm goroutine 3")
			time.Sleep(1 * time.Second)
			ch4 <- 1
		}
	}()

	go func() {
		for {
			<-ch4
			fmt.Println("I'm goroutine 4")
			time.Sleep(1 * time.Second)
			ch1 <- 1
		}
	}()

	select {}
}

func orderPrint2ByTimeDelay() {
	ch := make(chan struct{})
	for i := 1; i <= 4; i++ {
		go func(index int) {
			time.Sleep(time.Duration(index*10) * time.Millisecond)
			for {
				<-ch
				fmt.Printf("I am No %d Goroutine\n", index)
				time.Sleep(time.Second)
				ch <- struct{}{}
			}
		}(i)
	}
	ch <- struct{}{}
	time.Sleep(time.Minute)
}

func orderPrint3ByWrapper() {
	f := func(i int, input <-chan int, output chan<- int) {
		for {
			<-input
			fmt.Println(i)
			time.Sleep(time.Second)
			output <- 1
		}
	}
	c := [4]chan int{}
	for i := range []int{1, 2, 3, 4} {
		c[i] = make(chan int)
	}
	go f(1, c[3], c[0])
	go f(2, c[0], c[1])
	go f(3, c[1], c[2])
	go f(4, c[2], c[3])
	c[3] <- 1
	select {}
}

func orderPrint4ByIterator() {
	Print := func(recv chan int) {
		for {
			index, ok := <-recv
			if !ok {
				return
			}
			fmt.Println(index)
			index++
			if index == 5 {
				index = 1
			}

			time.Sleep(1 * time.Second)
			recv <- index
		}
	}
	c := make(chan int)
	go Print(c)
	go Print(c)
	go Print(c)
	go Print(c)
	c <- 1
	time.Sleep(20 * time.Second)
	close(c)
}

func orderPrint5ByIterator() {
	const chanNum int = 4
	chanArr := make([]chan int, chanNum)
	for i := 0; i < chanNum; i++ {
		ch := make(chan int, 1)
		chanArr[i] = ch
	}

	chanArr[0] <- 1
	for i := 0; i < chanNum; i++ {
		nextChanIdx := (i + 1) % chanNum
		go func(cur, next chan int, idx int) {
			for {
				<-cur
				time.Sleep(1 * time.Second)
				fmt.Printf("%d\n", idx+1)
				next <- 1
			}
		}(chanArr[i], chanArr[nextChanIdx], i)
	}
	select {}
}

func orderPrint1ByJob() {
	n := &jobs.NumChan{}
	n.JobNum(4)
}

func main() {
	orderPrint1ByJob()
}
