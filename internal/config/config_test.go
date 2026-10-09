package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLoadOrCreateAPIKey_CreatesThenReuses(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", apiKeyFileName)

	first, err := loadOrCreateAPIKey(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !isValidAPIKey(first) {
		t.Fatalf("generated key is not valid: len=%d", len(first))
	}

	second, err := loadOrCreateAPIKey(path)
	if err != nil {
		t.Fatalf("reuse: %v", err)
	}
	if first != second {
		t.Fatal("expected the persisted key to be reused")
	}

	if runtime.GOOS != "windows" {
		info, _ := os.Stat(path)
		if perm := info.Mode().Perm(); perm != 0o600 {
			t.Fatalf("expected mode 0600, got %o", perm)
		}
	}
}

func TestLoadOrCreateAPIKey_ReplacesWeakKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), apiKeyFileName)
	if err := os.WriteFile(path, []byte("a\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	key, err := loadOrCreateAPIKey(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key == "a" || !isValidAPIKey(key) {
		t.Fatal("weak stored key should have been replaced")
	}
	data, _ := os.ReadFile(path)
	if strings.TrimSpace(string(data)) != key {
		t.Fatal("replacement key was not persisted")
	}
}

func TestLoadOrCreateAPIKey_TightensLoosePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permissions not supported")
	}
	path := filepath.Join(t.TempDir(), apiKeyFileName)
	stored := strings.Repeat("ab", 32)
	if err := os.WriteFile(path, []byte(stored), 0o644); err != nil {
		t.Fatal(err)
	}

	key, err := loadOrCreateAPIKey(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != stored {
		t.Fatal("valid stored key should be kept")
	}
	info, _ := os.Stat(path)
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("expected mode tightened to 0600, got %o", perm)
	}
}

func TestLoadOrCreateAPIKey_RefusesSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need extra privileges on Windows")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "elsewhere")
	stored := strings.Repeat("cd", 32)
	if err := os.WriteFile(target, []byte(stored), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, apiKeyFileName)
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}

	key, err := loadOrCreateAPIKey(path)
	if err == nil {
		t.Fatal("expected an error for a symlinked key file")
	}
	if key == stored {
		t.Fatal("key must not be read through a symlink")
	}
	data, _ := os.ReadFile(target)
	if string(data) != stored {
		t.Fatal("symlink target must not be modified")
	}
}
