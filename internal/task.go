package worker

import (
	"time"
)

type Task interface {
	Timeout() *time.Duration
	Name() string // indicates the class of task
	Args() interface{}
}

// TODO: maybe task should have Serialize() and Deserialize() methods that are used to pull and push from the queue



// AddingTask, SubtractingTask
