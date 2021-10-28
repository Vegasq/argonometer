package main

import (
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"time"
)

var HASH = []byte("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$RCeht9aM1+cvYpPovpMschlqMUf1vWgxUMuCDS1rMSM")

const WORKERS = 100
const TASKS = 1000000

func worker(in chan int, result chan uint8) {
	for {
		i := <-in
		argon2.VerifyEncoded([]byte(fmt.Sprintf("plain%d", i)), HASH)
		result <- 1
	}
}

func dispatchTasks(tasks chan int) {
	for i := 0; i < TASKS; i++ {
		tasks <- i
	}
}

func startWorkers(tasks chan int, result chan uint8) {
	for i := 0; i < WORKERS; i++ {
		go worker(tasks, result)
	}
}

func sleep(start time.Time) {
	for time.Since(start) < time.Minute {
		// ... sleep
	}
}

func counter(c *int, result chan uint8) {
	for {
		<-result
		(*c) += 1
	}
}

func main() {
	var c int

	tasks := make(chan int, TASKS)
	result := make(chan uint8, TASKS)

	go counter(&c, result)
	start := time.Now()
	dispatchTasks(tasks)
	startWorkers(tasks, result)
	sleep(start)

	fmt.Printf("Argon2id per minute:\t%d\n", c)

}
