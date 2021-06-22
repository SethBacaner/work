test:
	go test ./...


worker:
	go build -o build/worker cmd/worker/main.go && ./build/worker

compile-taskgen:
	go build -o build/taskgen cmd/taskgen/main.go
