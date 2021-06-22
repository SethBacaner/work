package generate

import (
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type TemplateArgs struct {
	TaskName string
}

func GenerateTemplate(templateArgs TemplateArgs) {
	t := template.Must(template.New("task").Parse(taskTemplate))

	rc, wc, errCh := pipe.Commands(
		exec.Command("gofmt"),
		// exec.Command("goimports"),
	)
	go func() {
		select {
		case err, ok := <-errCh:
			if ok && err != nil {
				panic(err)
			}
		}
	}()
	t.Execute(wc, templateArgs)
	wc.Close()
	io.Copy(os.Stdout, rc)
}

var taskTemplate = `package gen

import (
	"context"
	"github.com/samuelbacaner/worker/internal"
)

// TODO: this needs to be generated more manually
type {{.TaskName}}Args struct {
	a int
	b int
}

type {{.TaskName}} interface {
	Invoke(ctx context.Context, args {{.TaskName}}Args) error
}

func Register{{.TaskName}}(manager worker.Manager, {{.TaskName}} {{.TaskName}}) {

	taskFn := func(ctx context.Context, args interface{}) error {
		{{.TaskName}}Args, ok := args.({{.TaskName}}Args)
		if !ok {
			// TODO: we need to blow up aggresively
		}

		return {{.TaskName}}.Invoke(ctx, {{.TaskName}}Args)
	}

	manager.RegisterTask("{{.TaskName}}", taskFn)
}

// TODO: need to generate EnqueuerCLient
// TODO: need to generate serializer and deserializer
`
