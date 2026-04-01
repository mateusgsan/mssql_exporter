package sql_exporter

import (
	"testing"

	"github.com/mateusgsan/mssql_exporter/config"
)

func TestSetColumnType(t *testing.T) {
	columnTypes := make(columnTypeMap)

	// Test setting a new column type
	err := setColumnType("test_ctx", "col1", columnTypeKey, columnTypes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if columnTypes["col1"] != columnTypeKey {
		t.Errorf("Expected columnTypeKey, got %v", columnTypes["col1"])
	}

	// Test setting the same column type again
	err = setColumnType("test_ctx", "col1", columnTypeKey, columnTypes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test setting a different column type for the same column (should fail)
	err = setColumnType("test_ctx", "col1", columnTypeValue, columnTypes)
	if err == nil {
		t.Errorf("Expected error for conflicting column types, got nil")
	}
}

func TestNewQuery(t *testing.T) {
	qc := &config.QueryConfig{
		Name: "test_query",
		Query: "SELECT 1",
	}

	mc := &config.MetricConfig{
		Name: "test_metric",
		KeyLabels: []string{"key_col"},
		Values: []string{"val_col"},
	}

	mf, err := NewMetricFamily("test_ctx", mc, nil)
	if err != nil {
		t.Fatalf("Failed to create MetricFamily: %v", err)
	}

	q, err := NewQuery("test_ctx", qc, mf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if q.config.Name != "test_query" {
		t.Errorf("Expected query name 'test_query', got %q", q.config.Name)
	}
	if len(q.metricFamilies) != 1 {
		t.Errorf("Expected 1 metric family, got %d", len(q.metricFamilies))
	}
	if q.columnTypes["key_col"] != columnTypeKey {
		t.Errorf("Expected key_col to be columnTypeKey")
	}
	if q.columnTypes["val_col"] != columnTypeValue {
		t.Errorf("Expected val_col to be columnTypeValue")
	}
}
