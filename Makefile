_build:
	go build -o stub ./cmd/stub/main.go

sidebar:
	go run ./cmd/sidecar/main.go
