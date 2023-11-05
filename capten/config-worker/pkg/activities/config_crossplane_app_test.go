package activities

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileValuesReplace(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "test.yaml")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.WriteString("https://github.com/intelops/capten-templates.git"); err != nil {
		t.Fatal(err)
	}
	file.Close()

	if err := replaceCaptenUrls(dir, "replaced"); err != nil {
		t.Fatal(err)
	}

	readBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(readBytes) != "replaced" {
		t.Fail()
	}
}
