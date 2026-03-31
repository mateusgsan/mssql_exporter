package sql_exporter

import (
	"testing"

	"github.com/free/sql_exporter/config"
	"gopkg.in/yaml.v2"
)

func TestNewJob(t *testing.T) {
	yamlData := []byte(`
job_name: test_job
collectors: []
static_configs:
  - targets:
      target1: "sqlserver://user:pass@localhost:1433"
      target2: "sqlserver://user:pass@dbserver2:1433"
    labels:
      env: prod
`)

	var jc config.JobConfig
	if err := yaml.Unmarshal(yamlData, &jc); err != nil {
		t.Fatalf("Failed to unmarshal JobConfig: %v", err)
	}

	gc := &config.GlobalConfig{}

	j, err := NewJob(&jc, gc)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if j == nil {
		t.Fatalf("Expected job, got nil")
	}

	targets := j.Targets()
	if len(targets) != 2 {
		t.Errorf("Expected 2 targets, got %d", len(targets))
	}
}

func TestNewJobWithCollectors(t *testing.T) {
	yamlData := []byte(`
job_name: test_job
collectors: [test_collector]
static_configs:
  - targets:
      target1: "sqlserver://user:pass@localhost:1433"
`)

	var jc config.JobConfig
	if err := yaml.Unmarshal(yamlData, &jc); err != nil {
		t.Fatalf("Failed to unmarshal JobConfig: %v", err)
	}

	gc := &config.GlobalConfig{}

	j, err := NewJob(&jc, gc)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	targets := j.Targets()
	if len(targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(targets))
	}
}
