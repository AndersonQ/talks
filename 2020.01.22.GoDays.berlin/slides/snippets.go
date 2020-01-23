package main

import (
	"fmt"
	"log"
	"strings"
)

// start_gogenerate OMIT
//go:generate obfuscate Credentials
// start_Credentials OMIT
type Credentials struct {
	username string
	password string
}

// end_Credentials OMIT
// end_gogenerate OMIT

// start_CredentialsString OMIT
func (c Credentials) String() string {
	return fmt.Sprintf("username: %s, password: %s",
		strings.Repeat("*", len(c.username)),
		strings.Repeat("*", len(c.password)))
}

// end_CredentialsString OMIT

func printDebug() {
	// start_CredentialsPrintDebugExec OMIT
	credentials := Credentials{username: "bob", password: "password123"}
	// start_CredentialsPrintDebug OMIT
	log.Printf("crecentials: %#v", credentials)
	// end_CredentialsPrintDebug OMIT
	// end_CredentialsPrintDebugExec OMIT
}
