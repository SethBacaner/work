package pool

import (
	"fmt"
	"sync"
)

type Runnable func()

type Pool interface {
	AvailableWorkers() int
	Execute(runnable Runnable)
}

func New(maxConcurrency int) Pool {
	return &poolImpl{
		maxConcurrency: maxConcurrency,
		mu:             &sync.Mutex{},
	}
}

type poolImpl struct {
	maxConcurrency int
	activeWorkers  int
	mu             *sync.Mutex
}

func (pi *poolImpl) AvailableWorkers() int {
	var available int
	pi.mu.Lock()
	available = pi.maxConcurrency - pi.activeWorkers
	pi.mu.Unlock()
	return available
}

func (pi *poolImpl) Execute(runnable Runnable) {
	go func() {
		fmt.Println("Starting task")
		pi.mu.Lock()
		pi.activeWorkers++
		pi.mu.Unlock()

		runnable()

		pi.mu.Lock()
		pi.activeWorkers--
		pi.mu.Unlock()
		fmt.Println("Finished task")
	}()
}
