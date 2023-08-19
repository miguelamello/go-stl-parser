package main

import (
	"os"
	"os/exec"
	"testing"
	"strings"
	"fmt"
)

func TestMain(t *testing.T) {

	expected := ""

	// Get current directory
	currentDir, _ := os.Getwd()

	// Test 1: Expecting the filepath as a command-line argument
	fmt.Println("Test 1: Expecting the filepath as a command-line argument")
	out, err := exec.Command(currentDir + "/bin/stlparser").Output()
	if err != nil { t.Error(err) }
	expected = "Please provide the filename to be analyzed"
	if strings.Contains(string(out), expected) == false {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}

	// Test 2: Provided filepath must be a valid file
	fmt.Println("Test 2: Provided filepath must be a valid file")
	out, err = exec.Command(currentDir+"/bin/stlparser", "invalidpath").Output()
	if err != nil { t.Error(err) }
	expected = "No such file or directory"
	if strings.Contains(string(out), expected) == false {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}

	// Test 3: Number of facets: 116 and Surface area: 15.54526855783991
	fmt.Println("Test 3: Moon.stl | Number of facets: 116 and Surface area: 15.54526855783991")
	out, err = exec.Command(currentDir+"/bin/stlparser", currentDir+"/shape/ascii/moon.stl").Output()
	if err != nil { t.Error(err) }
	expected = "Number of facets: 116"
	if strings.Contains(string(out), expected) == false {
		t.Errorf("Expected: %s, Got: %s", expected, string(out))
	} else {
		expected = "Surface area: 15.54526855783991"
		if strings.Contains(string(out), expected) == false {
			t.Errorf("Expected: %s, Got: %s", expected, string(out))
		}
	}
}

