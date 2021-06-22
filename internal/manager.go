package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TaskFunction func(ctx context.Context, args interface{}) error

type Manager interface {
	Start()
	Context() context.Context
	RegisterTask(taskName string, taskFunction TaskFunction) // eventually we want some argument here that allows us to do type checking via reflection e.g. https://github.com/grpc/grpc-go/blob/master/examples/helloworld/helloworld/helloworld_grpc.pb.go#L96
	GetRegisteredTaskFunction(taskName string) TaskFunction
}

type ManagerImpl struct {
	started  bool
	ctx      context.Context
	q        Queue
	pool     Pool
	registry map[string]TaskFunction
}

func New() *ManagerImpl {
	return &ManagerImpl{}
}

func (mi *ManagerImpl) Context() context.Context {
	return mi.ctx
}

func (mi *ManagerImpl) GetRegisteredTaskFunction(taskName string) TaskFunction {
	// TODO: we do need to check this incase the user never registered a task
	taskFn, _ := mi.registry[taskName]
	return taskFn
}

func (mi *ManagerImpl) Start() {
	if mi.started {
		// TODO: needs atomocity
		return
	}

	mi.started = true

	ctx, cancelFn := context.WithCancel(context.Background())
	mi.ctx = ctx

	go mi.poll(ctx)

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-sigCh

	fmt.Println()
	fmt.Println(sig)

	cancelFn()

	fmt.Println()
	fmt.Println("Exiting...")
}

func (mi *ManagerImpl) poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			if mi.pool.AvailableWorkers() <= 0 {
				continue
			}

			task, err := mi.q.GetTask(ctx)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if task == nil {
				continue
			}

			go mi.pool.DoTask(task)
		}
	}
}

func (mi *ManagerImpl) RegisterTask(taskName string, taskFunction TaskFunction) {
	// TODO: reflection on taskFunction to ensure that it is who it says it is
	mi.registry[taskName] = taskFunction
}
