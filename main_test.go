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

// mockModbusClient is a mock implementation of the modbus.Client interface.
type mockModbusClient struct {
	ReadHoldingRegistersFunc func(address, quantity uint16) ([]byte, error)
}

func (m *mockModbusClient) ReadHoldingRegisters(address, quantity uint16) ([]byte, error) {
	if m.ReadHoldingRegistersFunc != nil {
		return m.ReadHoldingRegistersFunc(address, quantity)
	}
	return nil, errors.New("ReadHoldingRegistersFunc not implemented")
}

// Implement other methods of the modbus.Client interface if needed.
func (m *mockModbusClient) ReadInputRegisters(address, quantity uint16) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) ReadCoils(address, quantity uint16) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) ReadDiscreteInputs(address, quantity uint16) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) WriteSingleCoil(address, value uint16) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) WriteMultipleCoils(address, quantity uint16, value []byte) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) WriteSingleRegister(address, value uint16) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) WriteMultipleRegisters(address, quantity uint16, value []byte) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) MaskWriteRegister(address, andMask, orMask uint16) ([]byte, error) {
	return nil, nil
}
func (m *mockModbusClient) ReadFIFOQueue(address uint16) ([]byte, error) {
	return nil, nil
}

func TestWebSocketHandler(t *testing.T) {
	t.Run("Successful read", func(t *testing.T) {
		mockClient := &mockModbusClient{
			ReadHoldingRegistersFunc: func(address, quantity uint16) ([]byte, error) {
				return []byte{4, 210}, nil // Corresponds to 1234
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
		mockClient := &mockModbusClient{
			ReadHoldingRegistersFunc: func(address, quantity uint16) ([]byte, error) {
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
	t.Run("Creates default config if not found", func(t *testing.T) {
		// Ensure no config file exists
		os.Remove("config.toml")

		cfg := loadConfig()
		assert.Equal(t, "localhost", cfg.ServerIP)
		assert.Equal(t, 5020, cfg.ServerPort)

		// Check that the file was created
		_, err := os.Stat("config.toml")
		assert.NoError(t, err)
		os.Remove("config.toml")
	})

	t.Run("Loads existing config", func(t *testing.T) {
		// Create a dummy config file
		os.Remove("config.toml")
		file, err := os.Create("config.toml")
		assert.NoError(t, err)

		customCfg := Config{
			ServerIP:     "192.168.1.100",
			ServerPort:   5020,
			StartAddress: 100,
			Quantity:     10,
			DelaySeconds: 5,
			WebUIPort:    9090,
		}
		err = toml.NewEncoder(file).Encode(customCfg)
		assert.NoError(t, err)
		file.Close()

		cfg := loadConfig()
		assert.Equal(t, "192.168.1.100", cfg.ServerIP)
		assert.Equal(t, 9090, cfg.WebUIPort)

		os.Remove("config.toml")
	})

	t.Run("Uses defaults for invalid config", func(t *testing.T) {
		// Create an invalid config file
		os.Remove("config.toml")
		invalidContent := `server_ip = "localhost"
server_port = "not-a-number"`
		err := os.WriteFile("config.toml", []byte(invalidContent), 0644)
		assert.NoError(t, err)

		cfg := loadConfig()
		// It should fall back to default values
		assert.Equal(t, 5020, cfg.ServerPort)

		os.Remove("config.toml")
	})
}
