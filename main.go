// Package main implements a web-based Modbus TCP browser.
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// webFS is the embedded file system for web assets (HTML, CSS, JS).
//
//go:embed all:web
var webFS embed.FS

// loadConfig loads the configuration from config.toml.
// If config.toml does not exist, it creates one with default values.
// If config.toml exists but is invalid, it logs a warning and uses default values.
func loadConfig() Config {
	// Default configuration
	cfg := Config{
		ServerIP:     "localhost",
		ServerPort:   5020,
		StartAddress: 0,
		Quantity:     2,
		DelaySeconds: 1,
		WebUIPort:    8080,
		SlaveID:      1,
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
	return cfg
}

// setupRouter creates and configures the Gin router.
// It sets up static file serving, HTML rendering, and WebSocket endpoint.
func setupRouter(cfg Config) *gin.Engine {
	router := gin.Default()

	// Use gzip middleware for asset compression.
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Serve static files from the embedded file system.
	staticFS, err := fs.Sub(webFS, "web/static")
	if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/static", http.FS(staticFS))

	// Serve the main index.html file.
	router.GET("/", func(c *gin.Context) {
		indexHTML, err := webFS.ReadFile("web/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	// Create a new Modbus client.
	client, err := NewModbusClient(&cfg)
	if err != nil {
		log.Fatalf("Failed to create Modbus client: %v", err)
	}
	err = client.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Modbus server: %v", err)
	}
	client.SetUnitId(cfg.SlaveID)

	// WebSocket endpoint for real-time Modbus data.
	wsHandler := &WsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for WebSocket connections
			},
		},
		cfg:          cfg,
		modbusClient: client,
	}
	router.GET("/ws", func(c *gin.Context) {
		wsHandler.ServeHTTP(c.Writer, c.Request)
	})

	return router
}

// main is the entry point of the application.
// It loads configuration, sets up the Gin router, and starts the web server.
func main() {
	// Set Gin to ReleaseMode if the "release" build tag is used
	if _, ok := os.LookupEnv("GIN_MODE"); !ok { // Check if GIN_MODE is not set via environment
		if gin.ReleaseMode == "release" { // This condition will be true if the "release" build tag is used
			gin.SetMode(gin.ReleaseMode)
		}
	}

	cfg := loadConfig()
	router := setupRouter(cfg)

	webAddress := fmt.Sprintf(":%d", cfg.WebUIPort)
	fmt.Printf("Starting web server on http://localhost%s\n", webAddress)
	if err := router.Run(webAddress); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
