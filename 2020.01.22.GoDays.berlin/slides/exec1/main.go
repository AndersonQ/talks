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

func (c Credentials) String() string {
	return fmt.Sprintf("username: %s, password: %s",
		strings.Repeat("*", len(c.username)),
		strings.Repeat("*", len(c.password)))
}

func main() {
	credentials := Credentials{username: "bob", password: "password123"}
	// start_CredentialsPrintDebugExec OMIT
	log.Printf("[DEBUG] crecentials: \n%#v", credentials)
	// end_CredentialsPrintDebugExec OMIT
}

// start_Stringer OMIT
// Stringer is implemented by any value that has a String method,
// which defines the ``native'' format for that value.
// The String method is used to print values passed as an operand // HL
// to any format that accepts a string or to an unformatted printer // HL
// such as Print.
type Stringer interface {
	String() string
}

// end_Stringer OMIT

// start_GoStringer OMIT
// GoStringer is implemented by any value that has a GoString method,
// which defines the Go syntax for that value.
// The GoString method is used to print values passed as an operand // HL
// to a %#v format. // HL
type GoStringer interface {
	GoString() string // HL
}

// end_GoStringer OMIT
