package sql_exporter

import (
	"testing"
)

func TestBoolToFloat64(t *testing.T) {
	if val := boolToFloat64(true); val != 1.0 {
		t.Errorf("Expected 1.0 for true, got %f", val)
	}
	if val := boolToFloat64(false); val != 0.0 {
		t.Errorf("Expected 0.0 for false, got %f", val)
	}
}
