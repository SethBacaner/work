package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sethbacaner/work/internal/pool"
	"github.com/sethbacaner/work/internal/queue"
)

type Manager interface {
	Start()
	Context() context.Context
	RegisterTask(taskName string, argsFactory ArgsFactory, taskFunction TaskFunction)
	Registry() *Registry
}

type ManagerImpl struct {
	started    bool
	ctx        context.Context
	cancelFunc context.CancelFunc
	q          queue.Queue
	pool       pool.Pool
	runner     Runner
	registry   *Registry
}

type ManagerOpt func(*ManagerImpl)

func WithQueue(queue queue.Queue) ManagerOpt {
	return func(m *ManagerImpl) {
		m.q = queue
	}
}

func NewManager(opts ...ManagerOpt) *ManagerImpl {

	ctx, cancelFunc := context.WithCancel(context.Background())

	mi := &ManagerImpl{
		started:    false,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		q:          nil,
		runner:     NewRunner(),
		pool:       pool.New(5),
		registry: &Registry{
			definitions: make(map[string]*TaskDefinition, 0),
			mu:          &sync.Mutex{},
		},
	}

	for _, opt := range opts {
		opt(mi)
	}

	return mi
}

func (mi *ManagerImpl) Context() context.Context {
	return mi.ctx
}

func (mi *ManagerImpl) Start() {
	if mi.started {
		// TODO: needs atomocity
		return
	}

	mi.started = true

	go mi.runner.Run(mi.ctx, mi.pool, mi.q, mi.registry)

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-sigCh

	fmt.Println()
	fmt.Println(sig)

	defer mi.cancelFunc()

	fmt.Println()
	fmt.Println("Exiting...")
}

func (mi *ManagerImpl) RegisterTask(taskName string, argsFactory ArgsFactory, taskFunction TaskFunction) {
	// TODO: reflection on taskFunction to ensure that it is who it says it is
	definition := TaskDefinition{
		Name:        taskName,
		TaskFn:      taskFunction,
		argsFactory: argsFactory,
	}
	mi.registry.Register(&definition)
}

func (mi *ManagerImpl) Registry() *Registry {
	return mi.registry
}
