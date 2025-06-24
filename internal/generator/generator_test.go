package generator

import (
	"go/format"
	"os"
	"path/filepath"
	"testing"
)

// define a sample interface for testing
// this interface is part of the test package so we can load it by package path "github.com/AbdouB/noopito/internal/testpkg".
type Example interface {
	DoSomething(a int) error
	Value() string
}

func TestGenerateNoOp(t *testing.T) {
	dir := t.TempDir()
	err := GenerateNoOp("github.com/AbdouB/noopito/internal/testpkg", "Example", dir)
	if err != nil {
		t.Fatalf("GenerateNoOp failed: %v", err)
	}
	outFile := filepath.Join(dir, "example_noop.go")
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	// attempt to format to ensure valid go code
	_, err = format.Source(data)
	if err != nil {
		t.Fatalf("output is not valid go: %v", err)
	}
}
