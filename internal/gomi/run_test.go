package gomi

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRun_TooManyArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run([]string{"a", "b"}, strings.NewReader(""), &stdout, &stderr)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var exit *ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("expected *ExitError, got %T (%v)", err, err)
	}
	if exit.Code != 64 {
		t.Errorf("exit code = %d, want 64", exit.Code)
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("err = %q, want it to contain %q", err.Error(), "usage")
	}
}
