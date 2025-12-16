# Modbus Browser

A command-line tool written in Go to browse and query Modbus TCP servers.

## Configuration

This tool can be configured using a `config.toml` file in the project root. If no `config.toml` is present, it will default to connecting to `localhost:502` and reading 2 registers starting from address 0.

Example `config.toml`:

```toml
server_ip = "localhost"
server_port = 5020
start_address = 4000
quantity = 2
```

## Usage

You can build and run the application using the provided `Makefile`.

*   **Build the application:**
    ```shell
    make build
    ```

*   **Run the application:**
    ```shell
    make run
    ```

The module path for this project is `github.com/wends155/modbusBrowser`.