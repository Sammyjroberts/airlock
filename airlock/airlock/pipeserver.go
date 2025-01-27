package airlock

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

type ProxyBase struct {
	SockPath string
}

func (p *ProxyBase) Call(methodName string, args ...interface{}) (interface{}, error) {
	conn, err := net.Dial("unix", p.SockPath)
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

// Helper function to convert various number types to int
func ToInt(v interface{}) (int, bool) {
	switch num := v.(type) {
	case int:
		return num, true
	case int8:
		return int(num), true
	case int16:
		return int(num), true
	case int32:
		return int(num), true
	case int64:
		return int(num), true
	case uint8:
		return int(num), true
	case uint16:
		return int(num), true
	case uint32:
		return int(num), true
	case uint64:
		return int(num), true
	case float32:
		return int(num), true
	case float64:
		return int(num), true
	default:
		return 0, false
	}
}
