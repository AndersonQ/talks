package main

import (
	"log"
	"os/exec"

	"github.com/AndersonQ/talks/2020.12.10.Golang.Live/pwd"
)

func main() {
	cmd := exec.Command("pwd")
	out, _ := cmd.Output()
	log.Printf("main pwd: %s", string(out))
	pwd.PWD()
}
