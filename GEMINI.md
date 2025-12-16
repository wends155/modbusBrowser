# Gemini Code Understanding

## Project Overview

This project, `modbusBrowser`, is a command-line tool developed in Go designed to interact with Modbus TCP servers. Its primary function is to connect to a specified Modbus TCP server and read a range of holding registers. The application's behavior, including the server's address, port, and the Modbus read parameters (start address and quantity of registers), is configurable via a `config.toml` file.

## Configuration

The application's settings are managed through a `config.toml` file. If this file is not found, the application uses default hardcoded values.

**Example `config.toml`:**

```toml
server_ip = "localhost"
server_port = 5020
start_address = 4000
quantity = 2
```

**Configurable Parameters:**

*   `server_ip`: Specifies the IP address or hostname of the Modbus TCP server to connect to.
*   `server_port`: Defines the port number on which the Modbus TCP server is listening.
*   `start_address`: The starting address of the holding registers from which the application will begin reading data.
*   `quantity`: The total number of holding registers to read, starting from `start_address`.

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

## Dependencies

The project relies on the following Go modules:

*   `github.com/goburrow/modbus`: This library provides the necessary functionalities for Modbus TCP communication, enabling the application to establish connections and interact with Modbus devices.
*   `github.com/BurntSushi/toml`: Used for parsing and loading configuration settings from the `config.toml` file, offering a structured and easy-to-manage way to configure the application.
