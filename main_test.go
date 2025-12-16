package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketHandler(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(websocketHandler))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Test for server info message
	var msg WebSocketMessage
	err = ws.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("Failed to read JSON from WebSocket: %v", err)
	}

	if msg.Type != "serverInfo" {
		t.Errorf("Expected message type 'serverInfo', got '%s'", msg.Type)
	}
}
