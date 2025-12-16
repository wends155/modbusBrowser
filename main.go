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

//go:embed all:web
var webFS embed.FS

// Global variables for configuration and WebSocket upgrader.
var (
	cfg      Config
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for WebSocket connections
		},
	}
)

// Config holds the application's configuration.
type Config struct {
	ServerIP     string `toml:"server_ip"`
	ServerPort   int    `toml:"server_port"`
	StartAddress uint16 `toml:"start_address"`
	Quantity     uint16 `toml:"quantity"`
	DelaySeconds int    `toml:"delay_seconds"`
	WebUIPort    int    `toml:"web_ui_port"`
}

// WebSocketMessage defines the structure for messages sent over the WebSocket.
type WebSocketMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// loadConfig loads the configuration from config.toml, creating the file with default values if it doesn't exist.
func loadConfig() {
	// Default configuration
	cfg = Config{
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
}

// setupRouter creates and configures the Gin router.
func setupRouter() *gin.Engine {
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
		websocketHandler(c.Writer, c.Request)
	})

	return router
}

func main() {
	// Set Gin to ReleaseMode if the "release" build tag is used
	if _, ok := os.LookupEnv("GIN_MODE"); !ok { // Check if GIN_MODE is not set via environment
		if gin.ReleaseMode == "release" { // This condition will be true if the "release" build tag is used
			gin.SetMode(gin.ReleaseMode)
		}
	}

	loadConfig()
	router := setupRouter()

	webAddress := fmt.Sprintf(":%d", cfg.WebUIPort)
	fmt.Printf("Starting web server on http://localhost%s\n", webAddress)
	router.Run(webAddress)
}
