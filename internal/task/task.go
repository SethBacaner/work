package task

import (
	"time"
)

type Task interface {
	Timeout() *time.Duration
	Name() string // indicates the class of task
	Args() interface{}
}
