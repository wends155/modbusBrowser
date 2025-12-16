build:
	@go build -o bin/modbusBrowser.exe

build-prod:
	@go build -ldflags="-s -w" -o bin/modbusBrowser.exe

run:
	@go run main.go

clean:
ifeq ($(OS),Windows_NT)
	@cmd /C "if exist bin rmdir /S /Q bin"
else
	@rm -rf bin
endif
