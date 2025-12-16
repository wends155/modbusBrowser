package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/goburrow/modbus"
)

type Config struct {
	ServerIP     string `toml:"server_ip"`
	ServerPort   int    `toml:"server_port"`
	StartAddress uint16 `toml:"start_address"`
	Quantity     uint16 `toml:"quantity"`
	DelaySeconds int    `toml:"delay_seconds"`
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func resetCursor() {
	fmt.Print("\033[H")
}

func main() {
	// Default configuration
	cfg := Config{
		ServerIP:     "localhost",
		ServerPort:   502,
		StartAddress: 0,
		Quantity:     2,
		DelaySeconds: 1,
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

	// Set up channel for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(cfg.DelaySeconds) * time.Second)
	defer ticker.Stop()

	// Clear the screen once before the loop
	clearScreen()
	fmt.Printf("Attempting to connect to Modbus TCP server at %s\n", address)
	fmt.Printf("Reading registers every %d second(s).\n", cfg.DelaySeconds)
	fmt.Println("Press Ctrl+C to exit.")

	for {
		select {
		case <-ticker.C:
			// Read holding registers based on config
			results, err := client.ReadHoldingRegisters(cfg.StartAddress, cfg.Quantity)
			if err != nil {
				fmt.Printf("Error reading holding registers: %v\n", err)
				continue
			}
			resetCursor()
			fmt.Printf("Reading holding registers from address %d (Quantity: %d):\033[K\n", cfg.StartAddress, cfg.Quantity)
			output := ""
			for i := 0; i < len(results); i += 2 {
				address := cfg.StartAddress + uint16(i/2)
				value := uint16(results[i])<<8 | uint16(results[i+1])
				output += fmt.Sprintf("%d:%d, ", address, value)
			}
			// Remove the trailing comma and space
			if len(output) > 2 {
				output = output[:len(output)-2]
			}
			fmt.Printf("%s\033[K\n", output)
		case <-sigChan:
			fmt.Println("\nExiting.")
			return
		}
	}
}
