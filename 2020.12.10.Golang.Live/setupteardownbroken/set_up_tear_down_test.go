package setupteardownbroken

import (
	"fmt"
	"os"
	"testing"
)

// start_TestMain OMIT
func TestMain(m *testing.M) {
	if testing.Verbose() { // HL
		fmt.Println("TestMain")
	}

	os.Exit(m.Run())
}

// end_TestMain OMIT
