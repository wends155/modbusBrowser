# Modbus Browser

A web-based tool written in Go to browse and query Modbus TCP servers.

## Overview

This application provides a web interface to continuously monitor a range of holding registers from a specified Modbus TCP server. The backend is a Go application using the Gin framework, and it streams data to the frontend in real-time using WebSockets. The entire web UI is embedded into the Go binary, making it a self-contained, single-executable application. The web UI displays the connected Modbus server's address and port, along with continuously updated register values and a timestamp of the last update.

## Configuration

This tool is configured using a `config.toml` file in the project root. If this file is not found, it will be automatically created with default values.

**Example `config.toml`:**

```toml
server_ip = "localhost"
server_port = 5020
start_address = 4000
quantity = 2
delay_seconds = 1
web_ui_port = 8080
```

**Configurable Parameters:**

*   `server_ip`: The IP address or hostname of the Modbus TCP server.
*   `server_port`: The port number of the Modbus TCP server.
*   `start_address`: The starting address of the holding registers to read.
*   `quantity`: The number of holding registers to read.
*   `delay_seconds`: The delay in seconds between each Modbus read.
*   `web_ui_port`: The port on which the web UI will be served.

## Usage

You can build and run the application using the provided `Makefile`.

*   **Build the application:**
    ```shell
    make build
    ```
    Compiles the Go source code and generates an executable in the `bin/` directory.

*   **Build the application for production:**
    ```shell
    make build-prod
    ```
    This command compiles the application with optimizations (`-ldflags="-s -w"`) and a `release` build tag. This ensures Gin runs in `ReleaseMode`, optimizing performance for production.

*   **Run the application:**
    ```shell
    make run
    ```
    This command uses `go run .` to compile and run the application, which will start the web server. You can then access the web UI by navigating to `http://localhost:<web_ui_port>` in your web browser.

## Testing

Unit tests are included to verify the application's functionality. You can run the tests using the following command:

```shell
go test ./...
```

## How it Works

*   The Go backend is organized into `main.go` for application setup and `handlers.go` for WebSocket and Modbus logic. It serves a simple HTML/CSS/JS frontend, including a favicon. The `index.html` file is read from the embedded filesystem and written directly to the HTTP response. Asset delivery is optimized using gzip compression.
*   The backend uses the Gin web framework. In production builds, Gin is configured to run in `ReleaseMode`.
*   The frontend establishes a WebSocket connection to the `/ws` endpoint on the backend.
*   The backend sends structured JSON messages over the WebSocket.
    *   Initially, a message containing the connected server's IP and port is sent.
    *   Subsequently, messages with continuously updated Modbus register data and a timestamp are sent.
*   The frontend receives and parses these JSON messages, then displays the server information, timestamp, and updates the Modbus data in the web UI.

Made by Wendell Saligan

The module path for this project is `github.com/wends155/modbusBrowser`.
