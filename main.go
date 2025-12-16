package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/goburrow/modbus"
)

type Config struct {
	ServerIP     string `toml:"server_ip"`
	ServerPort   int    `toml:"server_port"`
	StartAddress uint16 `toml:"start_address"`
	Quantity     uint16 `toml:"quantity"`
}

func main() {
	// Default configuration
	cfg := Config{
		ServerIP:     "localhost",
		ServerPort:   502,
		StartAddress: 0,
		Quantity:     2,
	}

	// Load configuration from toml file if it exists
	if _, err := os.Stat("config.toml"); err == nil {
		if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
			fmt.Printf("Error loading config.toml: %v\n", err)
		}
	}

	address := fmt.Sprintf("%s:%d", cfg.ServerIP, cfg.ServerPort)

	// Connect to a default Modbus TCP server if no config.toml
	// Address: localhost:502
	handler := modbus.NewTCPClientHandler(address)
	handler.SlaveId = 1
	client := modbus.NewClient(handler)

	fmt.Printf("Attempting to connect to Modbus TCP server at %s\n", address)

	// Example: Read holding registers based on config
	results, err := client.ReadHoldingRegisters(cfg.StartAddress, cfg.Quantity)
	if err != nil {
		fmt.Printf("Error reading holding registers: %v\n", err)
		return
	}

	fmt.Printf("Successfully read holding registers: %v\n", results)
}
