package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Config holds the application's configuration.
type Config struct {
	ServerIP     string `toml:"server_ip"`
	ServerPort   int    `toml:"server_port"`
	StartAddress uint16 `toml:"start_address"`
	Quantity     uint16 `toml:"quantity"`
	DelaySeconds int    `toml:"delay_seconds"`
	WebUIPort    int    `toml:"web_ui_port"`
	SlaveID      byte   `toml:"slave_id"`
}

// WebSocketMessage defines the structure for messages sent over the WebSocket.
type WebSocketMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// WsHandler is a handler for WebSocket connections.
// It manages the WebSocket upgrade, Modbus client, and data streaming.
type WsHandler struct {
	upgrader     websocket.Upgrader
	cfg          Config
	modbusClient ModbusClientInterface // Use the interface type
}

// readModbusData reads data from the Modbus server and formats it into a human-readable string.
func (h *WsHandler) readModbusData() (string, error) {
	results, err := h.modbusClient.ReadRegisters(h.cfg.StartAddress, h.cfg.Quantity)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, value := range results { // results is []uint16
		if i > 0 {
			sb.WriteString(", ")
		}
		address := h.cfg.StartAddress + uint16(i)
		sb.WriteString(fmt.Sprintf("%d:%d", address, value))
	}
	return sb.String(), nil
}

// ServeHTTP handles WebSocket connections.
// It upgrades the HTTP connection to a WebSocket connection,
// then continuously reads Modbus data and sends it to the client.
func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Send server info to the frontend.
	serverInfoMsg := WebSocketMessage{
		Type:    "serverInfo",
		Content: fmt.Sprintf("Server: %s:%d", h.cfg.ServerIP, h.cfg.ServerPort),
	}
	if err := conn.WriteJSON(serverInfoMsg); err != nil {
		log.Println("write server info:", err)
		return
	}

	ticker := time.NewTicker(time.Duration(h.cfg.DelaySeconds) * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return // client closed or read error
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			output, err := h.readModbusData()
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
}
