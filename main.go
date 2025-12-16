package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/goburrow/modbus"
	"github.com/gorilla/websocket"
)

//go:embed all:web
var webFS embed.FS

type Config struct {
	ServerIP     string `toml:"server_ip"`
	ServerPort   int    `toml:"server_port"`
	StartAddress uint16 `toml:"start_address"`
	Quantity     uint16 `toml:"quantity"`
	DelaySeconds int    `toml:"delay_seconds"`
	WebUIPort    int    `toml:"web_ui_port"`
}

type WebSocketMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for WebSocket connections
	},
}

func main() {
	// Set Gin to ReleaseMode if the "release" build tag is used
	if _, ok := os.LookupEnv("GIN_MODE"); !ok { // Check if GIN_MODE is not set via environment
		if gin.ReleaseMode == "release" { // This condition will be true if the "release" build tag is used
			gin.SetMode(gin.ReleaseMode)
		}
	}

	// Default configuration
	cfg := Config{
		ServerIP:     "localhost",
		ServerPort:   502,
		StartAddress: 0,
		Quantity:     2,
		DelaySeconds: 1,
		WebUIPort:    8080,
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

	router := gin.Default()

	// Use gzip middleware
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Serve static files
	staticFS, err := fs.Sub(webFS, "web/static")
		if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/static", http.FS(staticFS))

	// Serve index.html
	router.GET("/", func(c *gin.Context) {
		indexHTML, err := webFS.ReadFile("web/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		address := fmt.Sprintf("%s:%d", cfg.ServerIP, cfg.ServerPort)
		handler := modbus.NewTCPClientHandler(address)
		handler.SlaveId = 1
		client := modbus.NewClient(handler)

		// Send server info to frontend
		serverInfoMsg := WebSocketMessage{
			Type:    "serverInfo",
			Content: fmt.Sprintf("Server: %s:%d", cfg.ServerIP, cfg.ServerPort),
		}
		if err := conn.WriteJSON(serverInfoMsg); err != nil {
			log.Println("write server info:", err)
			return
		}

		ticker := time.NewTicker(time.Duration(cfg.DelaySeconds) * time.Second)
		defer ticker.Stop()

		for t := range ticker.C {
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

			output := ""
			for i := 0; i < len(results); i += 2 {
				address := cfg.StartAddress + uint16(i/2)
				value := uint16(results[i])<<8 | uint16(results[i+1])
				output += fmt.Sprintf("%d:%d, ", address, value)
			}
			if len(output) > 2 {
				output = output[:len(output)-2]
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
	})

	webAddress := fmt.Sprintf(":%d", cfg.WebUIPort)
	fmt.Printf("Starting web server on http://localhost%s\n", webAddress)
	router.Run(webAddress)
}
