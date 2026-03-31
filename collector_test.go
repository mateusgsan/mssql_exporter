package sql_exporter

import (
	"testing"
	"time"

	"github.com/free/sql_exporter/config"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

func TestNewCollector(t *testing.T) {
	yamlData := []byte(`
collector_name: test_collector
metrics:
  - metric_name: test_metric
    type: gauge
    help: Test help
    values: ["val1"]
    query: "SELECT 1"
`)

	var cc config.CollectorConfig
	err := yaml.Unmarshal(yamlData, &cc)
	if err != nil {
		t.Fatalf("Failed to unmarshal CollectorConfig: %v", err)
	}

	coll, err := NewCollector("test_ctx", &cc, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if coll == nil {
		t.Fatalf("Expected collector, got nil")
	}
}

func TestNewCachingCollector(t *testing.T) {
	rawColl := &collector{
		config: &config.CollectorConfig{
			MinInterval: model.Duration(10 * time.Second),
		},
		logContext: "test_ctx",
	}

	cc := newCachingCollector(rawColl)
	if cc == nil {
		t.Fatalf("Expected caching collector, got nil")
	}

	cachingColl, ok := cc.(*cachingCollector)
	if !ok {
		t.Fatalf("Expected *cachingCollector type")
	}

	if cachingColl.minInterval != 10*time.Second {
		t.Errorf("Expected minInterval 10s, got %v", cachingColl.minInterval)
	}
}
