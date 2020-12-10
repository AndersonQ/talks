package main

import (
	"fmt"
	"math/rand"
)

// start_naming OMIT

type MyType string

func (m MyType) MyMethod() { fmt.Println("MyMethod") }

func MyFunc() { fmt.Println("MyFunc") }

// end_naming OMIT

func Hello(name string) {
	fmt.Println("Hello,", name)
}

func Unordered() []int {
	if rand.Intn(10)%2 == 0 {
		return []int{1, 2}
	}
	return []int{2, 1}
}
