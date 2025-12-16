package main

import (
	"fmt"
	"log"
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

	configFilePath := "config.toml"

	// Check if config.toml exists, if not, create it with default values
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Printf("config.toml not found. Creating with default values.")
		file, err := os.Create(configFilePath)
		if err != nil {
			log.Fatalf("Failed to create config.toml: %v", err)
		}
		defer file.Close()

		if err := toml.NewEncoder(file).Encode(cfg); err != nil {
			log.Fatalf("Failed to write default config to config.toml: %v", err)
		}
	}

	// Load configuration from toml file
	if _, err := toml.DecodeFile(configFilePath, &cfg); err != nil {
		log.Printf("Error loading config.toml: %v. Using default configuration.", err)
	}

	address := fmt.Sprintf("%s:%d", cfg.ServerIP, cfg.ServerPort)

	// Connect to a default Modbus TCP server if no config.toml
	// Address: localhost:502
	handler := modbus.NewTCPClientHandler(address)
	handler.SlaveId = 1
	client := modbus.NewClient(handler)

	// Initial connection check
	if _, err := client.ReadHoldingRegisters(cfg.StartAddress, 1); err != nil {
		log.Fatalf("Failed to connect to Modbus server at %s: %v", address, err)
	}

	// Set up channel for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(cfg.DelaySeconds) * time.Second)
	defer ticker.Stop()

	// Clear the screen once before the loop
	clearScreen()
	fmt.Printf("Successfully connected to Modbus TCP server at %s\n", address)
	fmt.Printf("Reading registers every %d second(s).\n", cfg.DelaySeconds)
	fmt.Println("Press Ctrl+C to exit.")

	for {
		select {
		case <-ticker.C:
			// Read holding registers based on config
			results, err := client.ReadHoldingRegisters(cfg.StartAddress, cfg.Quantity)
			if err != nil {
				resetCursor()
				fmt.Printf("Error reading holding registers: %v\033[K\n", err)
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
