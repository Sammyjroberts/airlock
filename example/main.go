package main

import "fmt"

func main() {
	caller := NewMathCaller("unixmath")

	result, err := caller.Add(1, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Add result:", result)
	result, err = caller.Subtract(1, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Subtract result:", result)

	result, err = caller.Multiply(1, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Multiply result:", result)
	result, err = caller.Divide(1, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Divide result:", result)

}
