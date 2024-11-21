package test

import (
	"os"
	"testing"

	"github.com/printchard/go-json-parser/parser"
)

const testCasesPath = "cases/"
const validCasesPath = testCasesPath + "valid/"
const invalidCasesPath = testCasesPath + "invalid/"

func TestValid(t *testing.T) {
	dir, err := os.ReadDir(validCasesPath)
	if err != nil {
		t.Fatalf("Error reading test cases: %s", err)
	}

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		file, err := os.ReadFile(validCasesPath + entry.Name())
		if err != nil {
			t.Errorf("Error reading file %s: %s", entry.Name(), err)
		}

		p := parser.New(string(file))
		val, err := p.Parse()
		if err != nil {
			t.Errorf("Expected valid output, got: %s", err)
		}

		t.Log(val)
	}
}

func TestInvalid(t *testing.T) {
	dir, err := os.ReadDir(invalidCasesPath)
	if err != nil {
		t.Fatalf("Error reading test cases: %s", err)
	}

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		file, err := os.ReadFile(invalidCasesPath + entry.Name())
		if err != nil {
			t.Errorf("Error reading file %s: %s", entry.Name(), err)
		}

		p := parser.New(string(file))
		val, err := p.Parse()
		if err == nil {
			t.Errorf("Expected invalid output, got: %v", val)
		}

		t.Log(err)
	}
}
