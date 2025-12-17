package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// mockModbusClientAdapter is a mock implementation of the ModbusClientInterface.
type mockModbusClientAdapter struct {
	OpenFunc      func() error
	CloseFunc     func() error
	SetUnitIdFunc func(slaveID byte)

	ReadRegistersFunc func(address, quantity uint16) ([]uint16, error)
}

// implement the interface
func (m *mockModbusClientAdapter) Open() error {
	if m.OpenFunc == nil {
		return nil
	}
	return m.OpenFunc()
}

func (m *mockModbusClientAdapter) Close() error {
	if m.CloseFunc == nil {
		return nil
	}
	return m.CloseFunc()
}

func (m *mockModbusClientAdapter) SetUnitId(slaveID byte) {
	if m.SetUnitIdFunc == nil {
		return
	}
	m.SetUnitIdFunc(slaveID)
}

func (m *mockModbusClientAdapter) ReadRegisters(address, quantity uint16) ([]uint16, error) {
	if m.ReadRegistersFunc == nil {
		return nil, errors.New("ReadRegisters not implemented")
	}
	return m.ReadRegistersFunc(address, quantity)
}

func TestWebSocketHandler(t *testing.T) {
	t.Run("Successful read", func(t *testing.T) {
		mockClient := &mockModbusClientAdapter{
			ReadRegistersFunc: func(address, quantity uint16) ([]uint16, error) {
				return []uint16{1234}, nil
			},
		}
		wsHandler := &WsHandler{
			upgrader: websocket.Upgrader{},
			cfg: Config{
				StartAddress: 4000,
				Quantity:     1,
				DelaySeconds: 1,
			},
			modbusClient: mockClient,
		}

		server := httptest.NewServer(http.HandlerFunc(wsHandler.ServeHTTP))
		defer server.Close()
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		var serverInfoMsg WebSocketMessage
		err = ws.ReadJSON(&serverInfoMsg)
		assert.NoError(t, err)
		assert.Equal(t, "serverInfo", serverInfoMsg.Type)

		var modbusDataMsg WebSocketMessage
		err = ws.ReadJSON(&modbusDataMsg)
		assert.NoError(t, err)
		assert.Equal(t, "modbusData", modbusDataMsg.Type)
		assert.Equal(t, "4000:1234", modbusDataMsg.Content)
	})

	t.Run("Error read", func(t *testing.T) {
		mockClient := &mockModbusClientAdapter{
			ReadRegistersFunc: func(address, quantity uint16) ([]uint16, error) {
				return nil, errors.New("mock Modbus error")
			},
		}
		wsHandler := &WsHandler{
			upgrader: websocket.Upgrader{},
			cfg: Config{
				DelaySeconds: 1,
			},
			modbusClient: mockClient,
		}

		server := httptest.NewServer(http.HandlerFunc(wsHandler.ServeHTTP))
		defer server.Close()
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		var serverInfoMsg WebSocketMessage
		err = ws.ReadJSON(&serverInfoMsg)
		assert.NoError(t, err)
		assert.Equal(t, "serverInfo", serverInfoMsg.Type)

		var modbusErrorMsg WebSocketMessage
		err = ws.ReadJSON(&modbusErrorMsg)
		assert.NoError(t, err)
		assert.Equal(t, "modbusData", modbusErrorMsg.Type)
		assert.Contains(t, modbusErrorMsg.Content, "mock Modbus error")
	})
}

func TestLoadConfig(t *testing.T) {
	// Use a temporary working directory so tests don't modify the user's config.toml
	tmpDir, err := os.MkdirTemp("", "config_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	defer func() { _ = os.Chdir(cwd) }()

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	t.Run("Creates default config if not found", func(t *testing.T) {
		cfg := loadConfig()
		assert.Equal(t, "localhost", cfg.ServerIP)
		assert.Equal(t, 5020, cfg.ServerPort)
		assert.Equal(t, byte(1), cfg.SlaveID)

		_, err := os.Stat("config.toml")
		assert.NoError(t, err)
	})

	t.Run("Loads existing config", func(t *testing.T) {
		file, err := os.Create("config.toml")
		assert.NoError(t, err)

		customCfg := Config{
			ServerIP:     "192.168.1.100",
			ServerPort:   5020,
			StartAddress: 100,
			Quantity:     10,
			DelaySeconds: 5,
			WebUIPort:    9090,
			SlaveID:      2,
		}
		err = toml.NewEncoder(file).Encode(customCfg)
		assert.NoError(t, err)
		file.Close()

		cfg := loadConfig()
		assert.Equal(t, "192.168.1.100", cfg.ServerIP)
		assert.Equal(t, 9090, cfg.WebUIPort)
		assert.Equal(t, byte(2), cfg.SlaveID)
	})

	t.Run("Uses defaults for invalid config", func(t *testing.T) {
		invalidContent := `server_ip = "localhost"
server_port = "not-a-number"`
		err := os.WriteFile("config.toml", []byte(invalidContent), 0644)
		assert.NoError(t, err)

		cfg := loadConfig()
		assert.Equal(t, 5020, cfg.ServerPort)
		assert.Equal(t, byte(1), cfg.SlaveID) // Ensure SlaveID defaults correctly
	})
}
