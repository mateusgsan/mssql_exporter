package sql_exporter

import (
	"testing"

	"github.com/free/sql_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewTarget(t *testing.T) {
	gc := &config.GlobalConfig{}
	
	// Test with empty name (single target mode)
	target, err := NewTarget("test_ctx", "", "sqlserver://user:pass@localhost:1433", nil, nil, gc)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if target == nil {
		t.Fatalf("Expected target, got nil")
	}

	// Test with name
	labels := prometheus.Labels{"env": "prod"}
	target2, err := NewTarget("test_ctx", "test_target", "sqlserver://user:pass@localhost:1433", nil, labels, gc)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if target2 == nil {
		t.Fatalf("Expected target, got nil")
	}
}
