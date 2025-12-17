package main

import (
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

// ModbusClientInterface defines the operations needed from a Modbus client.

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
	config *Config
}

func IPtoURL(conf *Config) string {
	return fmt.Sprintf("tcp://%s:%d", conf.ServerIP, conf.ServerPort)
}

// NewModbusClient creates a new adapter for simonvetter/modbus.

func NewModbusClient(cfg *Config) (ModbusClientInterface, error) {

	client, err := modbus.NewClient(&modbus.ClientConfiguration{

		URL:     IPtoURL(cfg),
		Timeout: 1 * time.Second,
	})

	if err != nil {

		return nil, err

	}

	return &modbusClient{client: client, config: cfg}, nil

}

// Open implements the ModbusClientInterface Open method.

func (s *modbusClient) Open() error {

	return s.client.Open()

}

// Close implements the ModbusClientInterface Close method.

func (s *modbusClient) Close() error {

	return s.client.Close()

}

// SetUnitId implements the ModbusClientInterface SetUnitId method.

func (s *modbusClient) SetUnitId(slaveID byte) {

	s.client.SetUnitId(slaveID)

}

// ReadHoldingRegisters implements the ModbusClientInterface for simonvetter/modbus.

func (s *modbusClient) ReadHoldingRegisters(address, quantity uint16) ([]uint16, error) {

	return s.client.ReadRegisters(address, quantity, modbus.HOLDING_REGISTER)

}

func (s *modbusClient) ReadRegisters(addr uint16, quantity uint16) ([]uint16, error) {

	return s.client.ReadRegisters(addr, quantity, modbus.HOLDING_REGISTER)

}
