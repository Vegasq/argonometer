package main

import (
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"runtime/debug"
	"time"
)

var HASH = []byte("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$RCeht9aM1+cvYpPovpMschlqMUf1vWgxUMuCDS1rMSM")
const TASKS = 1000000

func worker(in chan int, result chan bool, exit chan bool) {
	for {
		select {
		case task := <- in:
			argon2.VerifyEncoded([]byte(fmt.Sprintf("plain%d", task)), HASH)
			result <- true
		case <- exit:
			return
		}
	}
}

func dispatchTasks(tasks chan int) {
	for i := 0; i < TASKS; i++ {
		tasks <- i
	}
}

func startWorkers(workers int, tasks chan int, result chan bool, exit chan bool) {
	//fmt.Printf("Start %d workers\n", workers)
	for i := 0; i < workers; i++ {
		go worker(tasks, result, exit)
	}
}

func sleep(start time.Time) {
	for time.Since(start) < time.Minute {
		// ... sleep
	}
}

func counter(c *int, result chan bool) {
	for {
		<-result
		(*c) += 1
	}
}

func benchmark(workers int) int {
	var c int

	tasks := make(chan int, TASKS)
	result := make(chan bool, TASKS)
	exit := make(chan bool)

	go counter(&c, result)
	start := time.Now()
	dispatchTasks(tasks)
	startWorkers(workers, tasks, result, exit)
	sleep(start)
	for i := 0; i < workers; i++ {
		exit <- true
	}

	return c
}


func main(){
	for i := 4; i <= 512; i = i*2 {
		hpm := benchmark(i)
		fmt.Printf("Workers: %d\tHashes: %d\n", i, hpm)
		debug.FreeOSMemory()
	}
}
