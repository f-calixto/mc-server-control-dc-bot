compile-arm64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build --ldflags="-w -s" -o ./bin/main