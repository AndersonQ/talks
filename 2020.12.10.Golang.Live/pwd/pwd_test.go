package pwd

import (
	"os/exec"
	"testing"
)

func TestPWD(t *testing.T) {
	cmd := exec.Command("pwd")
	out, _ := cmd.Output()
	t.Logf("pwd.TestPWD: %s", string(out))
	PWD()
}
