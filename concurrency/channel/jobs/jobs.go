package jobs

import (
	"fmt"
	"time"
)

type NumChan struct {
	jobs []*Job
}

func (n *NumChan) JobNum(m int) {
	for i := 1; i <= m; i++ {
		job := &Job{
			ID: i,
			Ch: make(chan struct{}, 1),
		}
		go job.run()
		n.jobs = append(n.jobs, job)
	}
	for {
		n.seq()
	}
}

func (n *NumChan) seq() {
	for _, j := range n.jobs {
		j.Ch <- struct{}{}
		time.Sleep(time.Second * 1)
	}
}

type Job struct {
	ID int
	Ch chan struct{}
}

func (j *Job) run() {
	// for {
	// 	select {
	// 	case <-j.Jobc:
	// 		fmt.Printf("id %d\n", j.ID)
	// 	}
	// }
	for range j.Ch {
		fmt.Printf("id %d\n", j.ID)
	}
}
