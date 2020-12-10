package setupteardown

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func globalSetup() {
	fmt.Println("Package level SetUp")
}
func globalTeardown() {
	fmt.Println("Package level TearDown")
}

func TestSetUpTearDown1(t *testing.T) {
	t.Log("SetUp/TearDown test 1")
}

func TestSetUpTearDown2(t *testing.T) {
	t.Log("SetUp/TearDown test 2")
	t.Log("SetUp/TearDown test 2")
}

// start_TestMain OMIT
func TestMain(m *testing.M) { // HL
	// call flag.Parse() here if TestMain uses flags

	// If TestMain depends on command-line flags, including those of the testing package,
	// it should call flag.Parse explicitly (https://godoc.org/testing)
	flag.Parse()

	if testing.Verbose() {
		fmt.Println("TestMain")
	}

	globalSetup() // HL
	exitStatus := m.Run()
	globalTeardown() // HL

	os.Exit(exitStatus) // HL
}

// end_TestMain OMIT
