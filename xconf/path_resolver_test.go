package xconf

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_ResolvePaths(t *testing.T) {
	type TConfig struct {
		UserDir string `xconf:"resolve"`
		RelPath string `xconf:"resolve"`
	}

	config := TConfig{
		UserDir: "~/files.txt",
		RelPath: "./files.txt",
	}

	err := ResolvePaths("/path/to/configdir/config.yaml", &config)
	if err != nil {
		t.Fatal(err)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home directory: %v", err)
	}

	if config.UserDir != filepath.Join(homedir, "files.txt") {
		t.Fatalf("expected %s, got %s", filepath.Join(homedir, "files.txt"), config.UserDir)
	}

	if config.RelPath != "/path/to/configdir/files.txt" {
		t.Fatalf("expected /path/to/configdir/files.txt, got %s", config.RelPath)
	}
}
