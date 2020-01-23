package main

import (
	"fmt"

	"github.com/AndersonQ/obfuscate/example/models"
)

func main() {
	s := models.Secret("a secret")
	long := models.Secret("a vere long secret")
	fmt.Println("the secret (%s):", s)
	fmt.Printf("the secret (%%v): %v\n", s)
	fmt.Printf("the secret (%%#v): %#v\n", s)

	fmt.Println("the long secret (%s):", long)
	fmt.Printf("the long secret (%%v): %v\n", long)
	fmt.Printf("the long secret (%%#v): %#v\n", long)

	as := models.AnotherSecret{
		Key1: "Key1",
		Key2: 2,
	}
	fmt.Println("another secret (%s):", as)
	fmt.Printf("another secret (%%v): %v\n", as)
	fmt.Printf("another secret (%%#v): %#v\n", as)
}
