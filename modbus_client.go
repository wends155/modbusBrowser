package main

import (
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

// ModbusClientInterface defines the operations needed from a Modbus client.
// This interface allows for easier migration between different Modbus client libraries.
type ModbusClientInterface interface {
	Open() error
	Close() error
	SetUnitId(slaveID byte)
	ReadRegisters(addr uint16, quantity uint16) ([]uint16, error)
	// Add other Modbus functions as needed
}

// modbusClient adapts github.com/simonvetter/modbus.ModbusClient to ModbusClientInterface.
type modbusClient struct {
	client *modbus.ModbusClient
	config *Config // Add config to access delaySeconds
}

// IPtoURL converts an IP address and port from the Config struct into a Modbus TCP URL string.
func IPtoURL(conf *Config) string {
	return fmt.Sprintf("tcp://%s:%d", conf.ServerIP, conf.ServerPort)
}

// NewModbusClient creates a new adapter for simonvetter/modbus.
// It initializes the underlying Modbus client with the provided configuration.
func NewModbusClient(cfg *Config) (ModbusClientInterface, error) {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     IPtoURL(cfg),
		Timeout: time.Duration(cfg.DelaySeconds) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &modbusClient{client: client, config: cfg}, nil
}

// Open implements the ModbusClientInterface Open method.
// It establishes a connection to the Modbus server.
func (s *modbusClient) Open() error {
	return s.client.Open()
}

// Close implements the ModbusClientInterface Close method.
// It closes the connection to the Modbus server.
func (s *modbusClient) Close() error {
	return s.client.Close()
}

// SetUnitId implements the ModbusClientInterface SetUnitId method.
// It sets the Modbus Unit ID (slave ID) for the client.
func (s *modbusClient) SetUnitId(slaveID byte) {
	s.client.SetUnitId(slaveID)
}

// ReadRegisters implements the ModbusClientInterface ReadRegisters method.
// It reads general registers from the Modbus server (in this context, it also reads holding registers).
func (s *modbusClient) ReadRegisters(addr uint16, quantity uint16) ([]uint16, error) {
	return s.client.ReadRegisters(addr, quantity, modbus.HOLDING_REGISTER)
}
