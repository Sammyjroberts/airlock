package main

import (
	"fmt"
	"net"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

type RPCMessage struct {
	Function string        `msgpack:"func"`
	Args     []interface{} `msgpack:"args"`
}

type RPCResponse struct {
	Result interface{} `msgpack:"result"`
	Error  string      `msgpack:"error,omitempty"`
}

type UnixCaller[T interface{}] struct {
	sockPath string
	proxy    T
}

type proxyBase struct {
	sockPath string
}

func (p *proxyBase) call(methodName string, args ...interface{}) (interface{}, error) {
	conn, err := net.Dial("unix", p.sockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket: %w", err)
	}
	defer conn.Close()

	msg := RPCMessage{
		Function: methodName,
		Args:     args,
	}

	encoder := msgpack.NewEncoder(conn)
	if err := encoder.Encode(msg); err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}

	var resp RPCResponse
	decoder := msgpack.NewDecoder(conn)
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	return resp.Result, nil
}

// Math-specific proxy
type MathProxy struct {
	*proxyBase
}

func (p *MathProxy) Add(a, b int) int {
	result, err := p.call("Add", a, b)
	fmt.Println("result:", result)
	if err != nil {
		fmt.Println("error:", err)
		return 0
	}
	if val, ok := result.(int8); ok {
		return int(val)
	} else {
		fmt.Println("invalid result type", reflect.TypeOf(result))
	}
	return 0
}

func (p *MathProxy) Subtract(a, b int) int {
	result, err := p.call("Subtract", a, b)

	if err != nil {
		return 0
	}
	if val, ok := result.(int); ok {
		return val
	}
	return 0
}

func (p *MathProxy) Multiply(a, b int) int {
	result, err := p.call("Multiply", a, b)
	if err != nil {
		return 0
	}
	if val, ok := result.(int); ok {
		return val
	}
	return 0
}

func (p *MathProxy) Divide(a, b int) (int, error) {
	result, err := p.call("Divide", a, b)
	if err != nil {
		return 0, err
	}
	if val, ok := result.(int); ok {
		return val, nil
	}
	return 0, fmt.Errorf("invalid result type")
}

func NewUnixCaller[T interface{}](sockPath string) (*UnixCaller[T], error) {
	base := &proxyBase{sockPath: sockPath}

	// For Math interface
	proxy := &MathProxy{base}

	caller := &UnixCaller[T]{
		sockPath: sockPath,
		proxy:    any(proxy).(T),
	}

	return caller, nil
}

func (c *UnixCaller[T]) Get() T {
	return c.proxy
}
