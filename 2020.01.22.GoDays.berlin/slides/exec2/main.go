package main

import (
	"fmt"
	"log"
	"strings"
)

type Credentials struct {
	username string
	password string
}

// start_GoStringer OMIT
func (c Credentials) String() string {
	return fmt.Sprintf("username: %s, password: %s",
		strings.Repeat("*", len(c.username)),
		strings.Repeat("*", len(c.password)))
}

func (c Credentials) GoString() string {
	return c.String()
}

// end_GoStringer OMIT

func main() {
	credentials := Credentials{username: "bob", password: "password123"}
	// start_CredentialsPrintDebugExec OMIT
	log.Printf("crecentials: \n%#v", credentials)
	// end_CredentialsPrintDebugExec OMIT
}
