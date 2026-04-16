package main

import (
	"os/exec"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("go", "build", "-o", t.TempDir()+"/vaultpull", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return t.TempDir() + "/vaultpull"
}

func TestMain_Version(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "vaultpull v") {
		t.Errorf("expected version string, got: %s", out)
	}
}

func TestMain_MissingConfig(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-config", "nonexistent.yaml")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config")
	}
	if !strings.Contains(string(out), "error loading config") {
		t.Errorf("expected config error message, got: %s", out)
	}
}

func TestVersion_Constant(t *testing.T) {
	if version == "" {
		t.Error("version constant should not be empty")
	}
	if !strings.HasPrefix(version, "0.") && !strings.HasPrefix(version, "1.") {
		t.Errorf("unexpected version format: %s", version)
	}
}
