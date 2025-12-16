package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/goburrow/modbus"
)

// websocketHandler handles WebSocket connections.
// It upgrades the HTTP connection to a WebSocket connection,
// then continuously reads Modbus data and sends it to the client.
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Create a new Modbus client.
	address := fmt.Sprintf("%s:%d", cfg.ServerIP, cfg.ServerPort)
	handler := modbus.NewTCPClientHandler(address)
	handler.SlaveId = 1
	client := modbus.NewClient(handler)

	// Send server info to the frontend upon successful connection.
	serverInfoMsg := WebSocketMessage{
		Type:    "serverInfo",
		Content: fmt.Sprintf("Server: %s:%d", cfg.ServerIP, cfg.ServerPort),
	}
	if err := conn.WriteJSON(serverInfoMsg); err != nil {
		log.Println("write server info:", err)
		return
	}

	// Create a ticker to read Modbus data at a regular interval.
	ticker := time.NewTicker(time.Duration(cfg.DelaySeconds) * time.Second)
	defer ticker.Stop()

	// Loop indefinitely, reading and sending Modbus data.
	for t := range ticker.C {
		// Read holding registers from the Modbus server.
		results, err := client.ReadHoldingRegisters(cfg.StartAddress, cfg.Quantity)
		if err != nil {
			log.Printf("Error reading holding registers: %v", err)
			errorMsg := WebSocketMessage{
				Type:      "modbusData",
				Content:   fmt.Sprintf("Error: %v", err),
				Timestamp: t.Format(time.RFC3339),
			}
			if err := conn.WriteJSON(errorMsg); err != nil {
				log.Println("write error data:", err)
				break
			}
			continue
		}

		// Format the Modbus data into a human-readable string.
		output := ""
		for i := 0; i < len(results); i += 2 {
			address := cfg.StartAddress + uint16(i/2)
			value := uint16(results[i])<<8 | uint16(results[i+1])
			output += fmt.Sprintf("%d:%d, ", address, value)
		}
		if len(output) > 2 {
			output = output[:len(output)-2]
		}

		// Send the Modbus data to the frontend.
		dataMsg := WebSocketMessage{
			Type:      "modbusData",
			Content:   output,
			Timestamp: t.Format(time.RFC3339),
		}
		if err := conn.WriteJSON(dataMsg); err != nil {
			log.Println("write modbus data:", err)
			break
		}
	}
}
