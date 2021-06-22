package main

import (
	"flag"
	"github.com/samuelbacaner/worker/internal/generate"
)

/*
Example usage:

make compile-taskgen
./build/taskgen -name=Adder > gen/Adder.go

*/

func main() {
	var templateArgs generate.TemplateArgs
	flag.StringVar(&templateArgs.TaskName, "name", "", "The task name for the task being generated")
	flag.Parse()

	generate.GenerateTemplate(templateArgs)
}
