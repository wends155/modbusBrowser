# Gemini Code Understanding

## Project Overview

This project, `modbusBrowser`, is a web-based tool developed in Go designed to interact with Modbus TCP servers. It provides a web interface to continuously monitor a range of holding registers from a specified Modbus TCP server. The application uses a Gin web server to serve the UI and a WebSocket connection to stream data in real-time. The UI now prominently displays the connected Modbus server's address and port, along with the continuously updated register values and a timestamp of the last update.

The application's behavior, including the server's address, port, read interval, and the Modbus read parameters (start address and quantity of registers), is configurable via a `config.toml` file.

## Architecture

*   **Backend:** A Go application using the Gin framework to serve a web UI and handle WebSocket connections.
*   **Frontend:** A simple HTML/CSS/JS single-page application that connects to the backend via a WebSocket to receive and display real-time Modbus data.
*   **Embedding:** The entire web UI (HTML, CSS, JS) is embedded into the Go binary using the `embed` package. This creates a self-contained, single-executable application.

## Configuration

The application's settings are managed through a `config.toml` file. If this file is not found, it will be automatically created with default values.

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

*   `server_ip`: Specifies the IP address or hostname of the Modbus TCP server to connect to.
*   `server_port`: Defines the port number on which the Modbus TCP server is listening.
*   `start_address`: The starting address of the holding registers from which the application will begin reading data.
*   `quantity`: The total number of holding registers to read, starting from `start_address`.
*   `delay_seconds`: The delay in seconds between each Modbus read operation.
*   `web_ui_port`: The port on which the web UI will be served.

## Building and Running

The project includes a `Makefile` to streamline common development and build tasks.

*   **Build the application:**
    ```shell
    make build
    ```
    This command compiles the Go source code and generates an executable named `modbusBrowser.exe` in the `bin/` directory.

*   **Build the application for production:**
    ```shell
    make build-prod
    ```
    This command compiles the application with optimizations (`-ldflags="-s -w"`) and the `release` build tag. When this build tag is present, Gin runs in `ReleaseMode`, which disables debug output and optimizes performance for production.

*   **Run the application:**
    ```shell
    make run
    ```
    This command starts the web server on the port specified by `web_ui_port` in the configuration.

## Development Conventions

*   **Gin Mode:** Gin runs in `DebugMode` by default, but it switches to `ReleaseMode` when the `release` build tag is used during compilation, improving performance and reducing logging for production deployments.
*   **Web UI:** The web UI is a single-page application served from the `/` route. The `index.html` file is read from the embedded filesystem and written to the HTTP response. The UI now displays the connected server's address and port, a timestamp of the last update, and the Modbus data.
*   **Asset Compression:** Gzip compression is enabled for static asset delivery to optimize performance.
*   **WebSocket Communication:**
    *   The frontend establishes a WebSocket connection to the `/ws` endpoint to receive real-time data.
    *   The backend sends structured JSON messages over the WebSocket.
    *   A message with `Type: "serverInfo"` and `Content: "Server: <ip>:<port>"` is sent upon connection to display server details.
    *   Messages with `Type: "modbusData"` contain the `address:value` formatted Modbus register data and a `Timestamp`.
*   **Error Handling:**
    *   **Configuration:** If `config.toml` is not found, it will be created with default values. If it is present but invalid, a warning is logged, and the application proceeds with the default configuration.
    *   **WebSocket Errors:** Errors during WebSocket communication are logged to the console.
    *   **Modbus Errors:** Errors during Modbus reads are sent over the WebSocket to be displayed in the UI.

## Dependencies

The project relies on the following Go modules:

*   `github.com/gin-gonic/gin`: A web framework for Go.
*   `github.com/gorilla/websocket`: A WebSocket implementation for Go.
*   `github.com/goburrow/modbus`: This library provides the necessary functionalities for Modbus TCP communication.
*   `github.com/BurntSushi/toml`: Used for parsing and loading configuration settings from the `config.toml` file.
*   `github.com/gin-contrib/gzip`: Gin middleware for Gzip compression.