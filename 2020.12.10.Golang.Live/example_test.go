package main

import "fmt"

// start_ExampleHello OMIT
func ExampleHello() {
	Hello("Jake")

	// Output:
	// Hello, Jake
}

// end_ExampleHello OMIT

// start_ExampleHello_broken OMIT
func ExampleHello_broken() {
	Hello("Jake")

	// Output:
	// Hello,  Jake
}

// end_ExampleHello_broken OMIT

// start_ExampleUnordered OMIT
func ExampleUnordered() {
	for _, n := range Unordered() {
		fmt.Println(n)
	}

	// Unordered output:
	// 1
	// 2
}

// end_ExampleUnordered OMIT

// start_namingexample OMIT
func ExampleMyFunc() {
	MyFunc()

	// Output:
	// MyFunc
}

func ExampleMyType() {
	m := MyType("Hello")
	fmt.Println(m)

	// Output:
	// Hello
}

func ExampleMyType_MyMethod() {
	MyType("Hello").MyMethod()

	// Output:
	// MyMethod
}

// end_namingexample OMIT
