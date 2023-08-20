_build:
	go build -o stub ./cmd/stub/main.go

sidecar:
	go run ./cmd/sidecar/main.go
