package main

import (
	"flag"
	"github.com/samuelbacaner/worker/internal/generate"
)

/*
Example usage:

make compile-taskgen
./build/taskgen -source=input/Adder.go -name=Adder > gen/Adder.go

*/

func main() {
	var taskName, source string

	flag.StringVar(&source, "source", "", "Source file with json serializable args struct named {name}Args")
	flag.StringVar(&taskName, "name", "", "The task name for the task being generated")
	flag.Parse()

	generate.GenerateTask(source, taskName)
}
