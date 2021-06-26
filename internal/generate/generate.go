package generate

import (
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type TemplateArgs struct {
	TaskName     string
	StructString string
}

func GenerateTask(source string, taskName string) {

	structString := getStructString(source)

	templateArgs := TemplateArgs{
		TaskName:     taskName,
		StructString: structString,
	}

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
	// t.Execute(os.Stdout, templateArgs)
	t.Execute(wc, templateArgs)
	wc.Close()
	io.Copy(os.Stdout, rc)

	// outputFile := "gen/" + taskName + ".go"
	// f, err := os.Create(outputFile)
	// if err != nil {
	// 	panic(err)
	// }

	// w := bufio.NewWriter(f)
	// io.Copy(w, rc)
}

var taskTemplate = `package gen

import (
	"context"
	"github.com/sethbacaner/worker/internal"
)

{{.StructString}}

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
