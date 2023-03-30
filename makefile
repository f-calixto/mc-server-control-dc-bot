compile-arm64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/main --ldflags="-w -s" cmd/main.go