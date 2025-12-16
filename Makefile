build:
	@go build -o bin/modbusBrowser.exe

build-prod:
	@go build -ldflags="-s -w" -o bin/modbusBrowser.exe

run:
	@go run main.go
