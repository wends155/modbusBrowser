# Modbus Browser

A command-line tool written in Go to browse and query Modbus TCP servers.

## Configuration

This tool can be configured using a `config.toml` file in the project root. If no `config.toml` is present, it will default to connecting to `localhost:502`, reading 2 registers starting from address 0, with a 1-second delay.

Example `config.toml`:

```toml
server_ip = "localhost"
server_port = 5020
start_address = 4000
quantity = 2
delay_seconds = 1
```

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
    Executes the application directly. This will continuously read Modbus registers at the interval specified by `delay_seconds` in the config. To stop the application, press `Ctrl+C` for a graceful exit.

*   **Clean the build directory:**
    ```shell
    make clean
    ```
    Removes the `bin/` directory and its contents.

The module path for this project is `github.com/wends155/modbusBrowser`.