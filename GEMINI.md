# Gemini Code Understanding

## Project Overview

This project, `modbusBrowser`, is a command-line tool developed in Go designed to interact with Modbus TCP servers. Its primary function is to connect to a specified Modbus TCP server and continuously read a range of holding registers at a configurable interval. The application's behavior, including the server's address, port, read interval, and the Modbus read parameters (start address and quantity of registers), is configurable via a `config.toml` file. The application also supports graceful shutdown upon receiving an interrupt signal (e.g., Ctrl+C) and provides a flicker-free console output by resetting the cursor position and clearing the line before each read operation. The output now presents the registers in a detailed `address:value` format.

## Configuration

The application's settings are managed through a `config.toml` file. If this file is not found, the application uses default hardcoded values.

**Example `config.toml`:**

```toml
server_ip = "localhost"
server_port = 5020
start_address = 4000
quantity = 2
delay_seconds = 1
```

**Configurable Parameters:**

*   `server_ip`: Specifies the IP address or hostname of the Modbus TCP server to connect to.
*   `server_port`: Defines the port number on which the Modbus TCP server is listening.
*   `start_address`: The starting address of the holding registers from which the application will begin reading data.
*   `quantity`: The total number of holding registers to read, starting from `start_address`.
*   `delay_seconds`: The delay in seconds between each Modbus read operation.

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
    This command also compiles the application but includes `ldflags="-s -w"` to strip debug information and symbol tables, resulting in a smaller executable suitable for production deployment. The output executable is `bin/modbusBrowser.exe`.

*   **Run the application:**
    ```shell
    make run
    ```
    This command uses `go run main.go` to execute the application directly from the source code. This is convenient for development and testing as it does not require an explicit build step beforehand.

## Development Conventions

*   **Continuous Reading:** The `main` function now includes a continuous loop that reads Modbus registers at a configurable interval. The output now provides a detailed breakdown of registers in `address:value` format, displaying the specific starting address and quantity of registers read.
*   **Graceful Shutdown:** The application registers a signal handler to gracefully exit when an interrupt signal (like `Ctrl+C`) is received, ensuring proper resource cleanup.
*   **Flicker-Free Screen Update:** To prevent screen flickering and mixed characters, the application uses a `resetCursor` function that prints an ANSI escape code (`\033[H`) to move the cursor to the top-left corner of the console before each read. After printing the data, it uses another ANSI escape code (`\033[K`) to clear the rest of the line. An initial screen clear is performed using a platform-aware `clearScreen` function.
*   **Error Handling:**
    *   **Configuration:** If `config.toml` is present but invalid, a warning is logged, and the application proceeds with the default configuration.
    *   **Initial Connection:** The application performs an initial Modbus read to verify the connection. If it fails, the program exits with a fatal error.
    *   **Read Errors:** During the continuous reading loop, any errors are printed to the console, and the application continues to the next read attempt. This makes the tool resilient to transient network issues.

## Dependencies

The project relies on the following Go modules:

*   `github.com/goburrow/modbus`: This library provides the necessary functionalities for Modbus TCP communication, enabling the application to establish connections and interact with Modbus devices.
*   `github.com/BurntSushi/toml`: Used for parsing and loading configuration settings from the `config.toml` file, offering a structured and easy-to-manage way to configure the application.
