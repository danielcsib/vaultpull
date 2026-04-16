package sync_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
)

func makeResult(vaultPath, envFile string, written int, err error) sync.Result {
	return sync.Result{
		Mapping: config.Mapping{VaultPath: vaultPath, EnvFile: envFile},
		Written: written,
		Err:     err,
	}
}

func TestPrintReport_AllSuccess(t *testing.T) {
	results := []sync.Result{
		makeResult("secret/app", ".env", 3, nil),
	}
	var buf strings.Builder
	sync.PrintReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "✓") {
		t.Errorf("expected success marker in output: %q", out)
	}
	if !strings.Contains(out, "3 keys") {
		t.Errorf("expected key count in output: %q", out)
	}
}

func TestPrintReport_WithFailure(t *testing.T) {
	results := []sync.Result{
		makeResult("secret/app", ".env", 0, errors.New("connection refused")),
	}
	var buf strings.Builder
	sync.PrintReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "✗") {
		t.Errorf("expected failure marker in output: %q", out)
	}
	if !strings.Contains(out, "connection refused") {
		t.Errorf("expected error message in output: %q", out)
	}
}

func TestHasErrors(t *testing.T) {
	if sync.HasErrors([]sync.Result{makeResult("a", "b", 1, nil)}) {
		t.Error("expected no errors")
	}
	if !sync.HasErrors([]sync.Result{makeResult("a", "b", 0, errors.New("oops"))}) {
		t.Error("expected errors detected")
	}
}
