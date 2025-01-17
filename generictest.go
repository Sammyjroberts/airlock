package main

import (
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

type UnixCallerTest[T interface{}] struct {
	sockPath string
	proxy    T
}

func (c *UnixCallerTest[T]) call(methodName string, args ...interface{}) (interface{}, error) {
	conn, err := net.Dial("unix", c.sockPath)
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
		return nil, errors.New(resp.Error)
	}

	return resp.Result, nil
}

func NewUnixCallerTest[T interface{}](sockPath string) (*UnixCallerTest[T], error) {
	caller := &UnixCallerTest[T]{
		sockPath: sockPath,
	}

	var t T
	tType := reflect.TypeOf(&t).Elem()

	// Create the proxy struct type
	proxyType := reflect.StructOf([]reflect.StructField{
		{
			Name: "Caller",
			Type: reflect.TypeOf(caller),
		},
	})

	// Create an instance of the proxy
	proxyValue := reflect.New(proxyType)
	proxyValue.Elem().Field(0).Set(reflect.ValueOf(caller))

	// Create a map to hold the method implementations
	methods := make(map[string]reflect.Value)

	// Implement each method
	for i := 0; i < tType.NumMethod(); i++ {
		method := tType.Method(i)

		// Create the method implementation
		methodType := method.Type
		methodImpl := reflect.MakeFunc(methodType, func(args []reflect.Value) []reflect.Value {
			// Skip the receiver in args
			callArgs := make([]interface{}, len(args)-1)
			for i := 1; i < len(args); i++ {
				callArgs[i-1] = args[i].Interface()
			}

			// Get the caller from the receiver's struct
			rcvr := args[0]
			caller := rcvr.Elem().Field(0).Interface().(*UnixCallerTest[T])

			// Make the RPC call
			result, err := caller.call(method.Name, callArgs...)

			// Handle the return values based on method signature
			if methodType.NumOut() == 1 {
				if result == nil {
					return []reflect.Value{reflect.Zero(methodType.Out(0))}
				}
				resultValue := reflect.ValueOf(result)
				if !resultValue.Type().ConvertibleTo(methodType.Out(0)) {
					return []reflect.Value{reflect.Zero(methodType.Out(0))}
				}
				return []reflect.Value{resultValue.Convert(methodType.Out(0))}
			} else if methodType.NumOut() == 2 {
				if err != nil {
					return []reflect.Value{
						reflect.Zero(methodType.Out(0)),
						reflect.ValueOf(err),
					}
				}
				resultValue := reflect.ValueOf(result)
				if !resultValue.Type().ConvertibleTo(methodType.Out(0)) {
					return []reflect.Value{
						reflect.Zero(methodType.Out(0)),
						reflect.ValueOf(fmt.Errorf("cannot convert result to required type")),
					}
				}
				return []reflect.Value{
					resultValue.Convert(methodType.Out(0)),
					reflect.Zero(methodType.Out(1)),
				}
			}
			panic("unexpected number of return values")
		})

		methods[method.Name] = methodImpl
	}

	// Create type implementing the interface
	interfaceType := reflect.StructOf([]reflect.StructField{
		{
			Name: "Impl",
			Type: proxyType,
		},
	})

	// Create instance of interface implementation
	implValue := reflect.New(interfaceType).Elem()
	implValue.Field(0).Set(proxyValue.Elem())

	// Convert to interface type
	impl := implValue.Addr().Interface()

	// Verify interface implementation
	if !reflect.TypeOf(impl).Implements(tType) {
		return nil, fmt.Errorf("failed to implement interface")
	}

	caller.proxy = impl.(T)
	return caller, nil
}

func (c *UnixCallerTest[T]) Get() T {
	return c.proxy
}
