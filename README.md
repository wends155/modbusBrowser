# Modbus Browser

A command-line tool written in Go to browse and query Modbus TCP servers.

## Configuration

This tool can be configured using a `config.toml` file in the project root. If this file is not found, it will be automatically created with default values.

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

## Usage

You can build and run the application using the provided `Makefile`.

*   **Build the application:**
    ```shell
    make build
    ```
    Compiles the Go source code and generates an executable in the `bin/` directory.

*   **Run the application:**
    ```shell
    make run
    ```
    Executes the application directly. This will continuously read Modbus registers at the interval specified by `delay_seconds` in the config. The output will now display each register in an `address:value` format (e.g., `4000:0, 4001:0`). To prevent flickering and mixed characters, the application resets the cursor to the top-left of the console and clears the line before each read. To stop the application, press `Ctrl+C` for a graceful exit.

*   **Clean the build directory:**
    ```shell
    make clean
    ```
    Removes the `bin/` directory and its contents.

## Error Handling

*   **Configuration:** If `config.toml` is not found, it will be created with default values. If it is present but contains errors, a warning will be displayed, and the application will fall back to its default settings.
*   **Initial Connection:** The application performs an initial connection test. If it cannot connect to the Modbus server on startup, it will exit with a fatal error.
*   **During Operation:** If a Modbus read fails, the error will be printed to the console, and the application will attempt to read again on the next cycle.

The module path for this project is `github.com/wends155/modbusBrowser`.
