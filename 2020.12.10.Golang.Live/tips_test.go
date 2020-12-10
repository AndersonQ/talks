package main

import (
	"errors"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func globalSetup()    {}
func globalTeardown() {}

// start_TestSample OMIT
func TestSample(t *testing.T) {
	// your test
}

// end_TestSample OMIT

// start_TestMain OMIT
func TestMain(m *testing.M) { // HL
	// call flag.Parse() here if TestMain uses flags

	globalSetup()
	exitCode := m.Run() // HL
	globalTeardown()

	os.Exit(exitCode) // HL
}

// end_TestMain OMIT

// start_TestSimple OMIT
func TestSimple(t *testing.T) {
	expected := 2
	actual := 3
	if actual != expected {
		t.Errorf("expected: %d, actual: %d", expected, actual) // HL
	}
}

// end_TestSimple OMIT

// start_TestSimpleFatal OMIT
func TestSimpleFatal(t *testing.T) {
	someSetup := errors.New("setup failed")
	expected, actual := 3, 2

	if someSetup != nil {
		t.Fatalf("set up failed, aborting test: %v", someSetup) // HL
	}
	// Not executed // HL
	if actual != expected {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}
}

// end_TestSimpleFatal OMIT

// start_TestVerbose OMIT
func TestVerbose(t *testing.T) {
	t.Log("only printed in verbose mode") // HL
	log.Println("log.Println: always printed")

	if testing.Verbose() { // HL
		log.Println("some verbose, but really useful info")
	}
}

// end_TestVerbose OMIT

// start_TestShort OMIT
func TestShort(t *testing.T) {
	// Can't be quicker
}

func TestShortNotSoShort(t *testing.T) {
	if testing.Short() { // HL
		t.Skip("Skip in short mode") // needs verbose flag to show this message // HL
	}
	time.Sleep(3 * time.Second)
}

// end_TestShort OMIT

// start_helperFunction OMIT
func TestWithHelper(t *testing.T) {
	helperFunction(t, "a parameter")
	// ...
}

func helperFunction(t *testing.T, param string) string {
	// start_helperFunction_line OMIT
	t.Fatal("helperFunction setup failed") // HL
	// end_helperFunction_line OMIT
	return param
}

// end_helperFunction OMIT

// start_betterHelperFunction OMIT
func TestWithBetterHelper(t *testing.T) {
	betterHelperFunction(t, "file.json")
}

func betterHelperFunction(t *testing.T, fileName string) string {
	t.Helper() // HL
	t.Fatal("betterHelperFunction setup failed")
	return ""
}

// end_betterHelperFunction OMIT

// start_TestTable OMIT
func TestTableSimple(t *testing.T) {
	// start_TestTable_tcs OMIT
	tcs := []struct {
		name string
		val  float64
		want float64
	}{
		{name: "the positive", val: 42, want: 42},  // HLtcs
		{name: "the negative", val: -42, want: 42}, // HLtcs
		{name: "zero", val: 0, want: 0},            // HLtcs
	}
	// end_TestTable_tcs OMIT
	// start_TestTable_for OMIT
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) { // HLrun
			got := math.Abs(tc.val)
			if tc.want != got {
				t.Errorf("want: %f, got: %f", tc.want, got)
			}
		})
	}
	// end_TestTable_for OMIT
}

// end_TestTable OMIT

// start_TestTableSlow OMIT
func TestTableSlow(t *testing.T) {
	tcs := []struct {
		name  string
		sleep time.Duration
	}{
		{name: "1s", sleep: 1 * time.Second},
		{name: "2s", sleep: 2 * time.Second},
		{name: "3s", sleep: 3 * time.Second},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			time.Sleep(tc.sleep)
		})
	}
}

// end_TestTableSlow OMIT

// start_TestTableParallel OMIT
func TestTableParallel(t *testing.T) {
	tcs := []struct {
		name  string
		sleep time.Duration
	}{
		{name: "1s", sleep: 1 * time.Second},
		{name: "2s", sleep: 2 * time.Second},
		{name: "3s", sleep: 3 * time.Second},
	}
	// start_for_TestTableParallel OMIT
	for _, tc := range tcs { // HLsharing
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()               // HL
			t.Log("running ", tc.name) // HLsharing
			time.Sleep(tc.sleep)
		})
	}
	// end_for_TestTableParallel OMIT
}

// end_TestTableParallel OMIT

// start_TestTableParallelFixed OMIT
func TestTableParallelFixed(t *testing.T) {
	tcs := []struct {
		name  string
		sleep time.Duration
	}{
		{name: "1s", sleep: 1 * time.Second},
		{name: "2s", sleep: 2 * time.Second},
		{name: "3s", sleep: 3 * time.Second},
	}
	// start_for_TestTableParallelFixed OMIT
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc // rebidding tc so the goroutines will not share it // HL
			t.Parallel()
			t.Log("running ", tc.name)
			time.Sleep(tc.sleep)
		})
	}
	// end_for_TestTableParallelFixed OMIT
}

// end_TestTableParallelFixed OMIT

// start_TestTestdata OMIT
func TestTestdata(t *testing.T) {
	expected := "Hello, golpher!"

	path := filepath.Join("testdata", "hello") // HL
	file, _ := ioutil.ReadFile(path)

	str := string(file)
	if str != expected {
		t.Fatalf("expected: %s, actual: %s", expected, str)
	}
}

// end_TestTestdata OMIT
