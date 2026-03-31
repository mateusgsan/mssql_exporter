package errors

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := "test_context"
	msg := "test error message"
	err := New(ctx, msg)

	if err.Context() != ctx {
		t.Errorf("Expected context %q, got %q", ctx, err.Context())
	}
	if err.RawError() != msg {
		t.Errorf("Expected raw error %q, got %q", msg, err.RawError())
	}
	expectedErrStr := fmt.Sprintf("[%s] %s", ctx, msg)
	if err.Error() != expectedErrStr {
		t.Errorf("Expected error string %q, got %q", expectedErrStr, err.Error())
	}
}

func TestErrorf(t *testing.T) {
	ctx := "test_context"
	err := Errorf(ctx, "error %d", 42)

	if err.Context() != ctx {
		t.Errorf("Expected context %q, got %q", ctx, err.Context())
	}
	if err.RawError() != "error 42" {
		t.Errorf("Expected raw error %q, got %q", "error 42", err.RawError())
	}
}

func TestWrap(t *testing.T) {
	// Wrap nil
	if Wrap("ctx", nil) != nil {
		t.Error("Expected Wrap(nil) to return nil")
	}

	// Wrap standard error
	stdErr := fmt.Errorf("standard error")
	wrapped := Wrap("ctx", stdErr)
	if wrapped.Context() != "ctx" {
		t.Errorf("Expected context 'ctx', got %q", wrapped.Context())
	}
	if wrapped.RawError() != "standard error" {
		t.Errorf("Expected raw error 'standard error', got %q", wrapped.RawError())
	}

	// Wrap WithContext error
	alreadyWrapped := New("old_ctx", "old error")
	wrappedAgain := Wrap("new_ctx", alreadyWrapped)
	if wrappedAgain.Context() != "old_ctx" {
		t.Errorf("Expected context 'old_ctx', got %q", wrappedAgain.Context())
	}
}

func TestWrapf(t *testing.T) {
	// Wrapf nil
	if Wrapf("ctx", nil, "prefix") != nil {
		t.Error("Expected Wrapf(nil) to return nil")
	}

	// Wrapf standard error
	stdErr := fmt.Errorf("standard error")
	wrapped := Wrapf("ctx", stdErr, "prefix %d", 1)
	if wrapped.Context() != "ctx" {
		t.Errorf("Expected context 'ctx', got %q", wrapped.Context())
	}
	if wrapped.RawError() != "prefix 1: standard error" {
		t.Errorf("Expected raw error 'prefix 1: standard error', got %q", wrapped.RawError())
	}

	// Wrapf WithContext error
	alreadyWrapped := New("old_ctx", "old error")
	wrappedAgain := Wrapf("new_ctx", alreadyWrapped, "prefix")
	if wrappedAgain.Context() != "old_ctx" {
		t.Errorf("Expected context 'old_ctx', got %q", wrappedAgain.Context())
	}
	if wrappedAgain.RawError() != "prefix: old error" {
		t.Errorf("Expected raw error 'prefix: old error', got %q", wrappedAgain.RawError())
	}
}
