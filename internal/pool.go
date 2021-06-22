package worker

import "sync"

type Pool interface {
	AvailableWorkers() int
	DoTask(task *Task)
}

type PoolImpl struct {
	runner         Runner
	maxConcurrency int
	activeWorkers  int
	mu             *sync.Mutex
}

func (pi *PoolImpl) AvailableWorkers() int {
	var available int
	pi.mu.Lock()
	available = pi.maxConcurrency - pi.activeWorkers
	pi.mu.Unlock()
	return available
}

func (pi *PoolImpl) DoTask(task Task) {
	pi.mu.Lock()
	pi.activeWorkers++
	pi.mu.Unlock()

	pi.runner.DoTask(task)

	pi.mu.Lock()
	pi.activeWorkers--
	pi.mu.Unlock()
}
