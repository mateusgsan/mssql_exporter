package sql_exporter

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewExporter(t *testing.T) {
	// Create a temporary config file
	configContent := []byte(`
global:
  scrape_timeout: 10s
  min_interval: 0s
  max_connections: 3
  max_idle_connections: 3
target:
  data_source_name: "sqlserver://user:pass@localhost:1433"
  collectors: ["test_collector"]
collectors:
  - collector_name: test_collector
    metrics:
      - metric_name: test_metric
        type: gauge
        help: Test help
        values: ["val1"]
        query: "SELECT 1 AS val1"
`)
	tmpfile, err := ioutil.TempFile("", "config.*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test NewExporter
	exporter, err := NewExporter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if exporter == nil {
		t.Fatalf("Expected exporter, got nil")
	}

	// Test Config()
	cfg := exporter.Config()
	if cfg == nil {
		t.Fatalf("Expected config, got nil")
	}
	if cfg.Target.DSN != "sqlserver://user:pass@localhost:1433" {
		t.Errorf("Expected DSN 'sqlserver://user:pass@localhost:1433', got %q", cfg.Target.DSN)
	}

	// Test WithContext()
	ctx := context.WithValue(context.Background(), "key", "value")
	exporterWithCtx := exporter.WithContext(ctx)
	if exporterWithCtx == nil {
		t.Fatalf("Expected exporter with context, got nil")
	}
	
	// We can't easily check the context inside the unexported struct, but we can verify it returns a valid Exporter
	if exporterWithCtx.Config() != cfg {
		t.Errorf("Expected config to be the same")
	}
}
