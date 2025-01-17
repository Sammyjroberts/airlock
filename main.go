package main

import "fmt"

func main() {
	caller, err := NewUnixCaller[Math]("math.sock")
	if err != nil {
		panic(err)
	}
	sum := caller.proxy.Add(1, 2)
	fmt.Println("1 + 2 =", sum)
}
