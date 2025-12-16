build:
	@go build -o bin/modbusBrowser.exe

build-prod:
	@go build -tags release -ldflags="-s -w" -o bin/modbusBrowser.exe

run:
	# This will start the web server on http://localhost:8080
	@go run main.go
