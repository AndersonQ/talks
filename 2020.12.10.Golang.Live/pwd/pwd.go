package pwd

import (
	"log"
	"os/exec"
)

func PWD() {
	cmd := exec.Command("pwd")
	out, _ := cmd.Output()
	log.Printf("pwd.PWD: %s", string(out))
}
