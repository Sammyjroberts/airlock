package main

import (
	"fmt"
	"net"

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
		return nil, fmt.Errorf("remote error: %s", resp.Error)
	}

	return resp.Result, nil
}
