package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

func TestGlobalConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
min_interval: 5s
scrape_timeout: 15s
scrape_timeout_offset: 1s
max_connections: 10
max_idle_connections: 5
`)

	var g GlobalConfig
	err := yaml.Unmarshal(yamlData, &g)
	if err != nil {
		t.Fatalf("Failed to unmarshal GlobalConfig: %v", err)
	}

	if g.MinInterval != model.Duration(5*time.Second) {
		t.Errorf("Expected MinInterval 5s, got %v", g.MinInterval)
	}
	if g.ScrapeTimeout != model.Duration(15*time.Second) {
		t.Errorf("Expected ScrapeTimeout 15s, got %v", g.ScrapeTimeout)
	}
	if g.TimeoutOffset != model.Duration(1*time.Second) {
		t.Errorf("Expected TimeoutOffset 1s, got %v", g.TimeoutOffset)
	}
	if g.MaxConns != 10 {
		t.Errorf("Expected MaxConns 10, got %d", g.MaxConns)
	}
	if g.MaxIdleConns != 5 {
		t.Errorf("Expected MaxIdleConns 5, got %d", g.MaxIdleConns)
	}
}

func TestGlobalConfigDefaults(t *testing.T) {
	yamlData := []byte(`{}`)

	var g GlobalConfig
	err := yaml.Unmarshal(yamlData, &g)
	if err != nil {
		t.Fatalf("Failed to unmarshal GlobalConfig: %v", err)
	}

	if g.MinInterval != model.Duration(0) {
		t.Errorf("Expected default MinInterval 0, got %v", g.MinInterval)
	}
	if g.ScrapeTimeout != model.Duration(10*time.Second) {
		t.Errorf("Expected default ScrapeTimeout 10s, got %v", g.ScrapeTimeout)
	}
	if g.TimeoutOffset != model.Duration(500*time.Millisecond) {
		t.Errorf("Expected default TimeoutOffset 500ms, got %v", g.TimeoutOffset)
	}
	if g.MaxConns != 3 {
		t.Errorf("Expected default MaxConns 3, got %d", g.MaxConns)
	}
	if g.MaxIdleConns != 3 {
		t.Errorf("Expected default MaxIdleConns 3, got %d", g.MaxIdleConns)
	}
}

func TestTargetConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
data_source_name: "sqlserver://user:pass@localhost:1433"
collectors: ["collector1", "collector2"]
`)

	var target TargetConfig
	err := yaml.Unmarshal(yamlData, &target)
	if err != nil {
		t.Fatalf("Failed to unmarshal TargetConfig: %v", err)
	}

	if target.DSN != "sqlserver://user:pass@localhost:1433" {
		t.Errorf("Expected DSN 'sqlserver://user:pass@localhost:1433', got %q", target.DSN)
	}
	if len(target.CollectorRefs) != 2 || target.CollectorRefs[0] != "collector1" || target.CollectorRefs[1] != "collector2" {
		t.Errorf("Expected collectors [collector1, collector2], got %v", target.CollectorRefs)
	}
}

func TestTargetConfigMissingDSN(t *testing.T) {
	yamlData := []byte(`
collectors: ["collector1"]
`)

	var target TargetConfig
	err := yaml.Unmarshal(yamlData, &target)
	if err == nil {
		t.Error("Expected error for missing DSN, got nil")
	}
}

func TestMetricConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
metric_name: test_metric
type: gauge
help: Test help
key_labels: ["label1"]
values: ["val1"]
query: "SELECT 1"
`)

	var m MetricConfig
	err := yaml.Unmarshal(yamlData, &m)
	if err != nil {
		t.Fatalf("Failed to unmarshal MetricConfig: %v", err)
	}

	if m.Name != "test_metric" {
		t.Errorf("Expected Name 'test_metric', got %q", m.Name)
	}
	if m.TypeString != "gauge" {
		t.Errorf("Expected TypeString 'gauge', got %q", m.TypeString)
	}
	if m.Help != "Test help" {
		t.Errorf("Expected Help 'Test help', got %q", m.Help)
	}
	if len(m.KeyLabels) != 1 || m.KeyLabels[0] != "label1" {
		t.Errorf("Expected KeyLabels [label1], got %v", m.KeyLabels)
	}
	if len(m.Values) != 1 || m.Values[0] != "val1" {
		t.Errorf("Expected Values [val1], got %v", m.Values)
	}
	if m.QueryLiteral != "SELECT 1" {
		t.Errorf("Expected QueryLiteral 'SELECT 1', got %q", m.QueryLiteral)
	}
}

func TestQueryConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
query_name: test_query
query: "SELECT 1"
`)

	var q QueryConfig
	err := yaml.Unmarshal(yamlData, &q)
	if err != nil {
		t.Fatalf("Failed to unmarshal QueryConfig: %v", err)
	}

	if q.Name != "test_query" {
		t.Errorf("Expected Name 'test_query', got %q", q.Name)
	}
	if q.Query != "SELECT 1" {
		t.Errorf("Expected Query 'SELECT 1', got %q", q.Query)
	}
}

func TestCollectorConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
collector_name: test_collector
min_interval: 10s
metrics:
  - metric_name: test_metric
    type: gauge
    help: Test help
    values: ["val1"]
    query: "SELECT 1"
queries:
  - query_name: test_query
    query: "SELECT 1"
`)

	var c CollectorConfig
	err := yaml.Unmarshal(yamlData, &c)
	if err != nil {
		t.Fatalf("Failed to unmarshal CollectorConfig: %v", err)
	}

	if c.Name != "test_collector" {
		t.Errorf("Expected Name 'test_collector', got %q", c.Name)
	}
	if c.MinInterval != model.Duration(10*time.Second) {
		t.Errorf("Expected MinInterval 10s, got %v", c.MinInterval)
	}
	if len(c.Metrics) != 1 {
		t.Errorf("Expected 1 metric, got %d", len(c.Metrics))
	}
	if len(c.Queries) != 1 {
		t.Errorf("Expected 1 query, got %d", len(c.Queries))
	}
}

func TestJobConfigUnmarshal(t *testing.T) {
	yamlData := []byte(`
job_name: test_job
collectors: ["collector1"]
static_configs:
  - targets:
      target1: "sqlserver://user:pass@localhost:1433"
    labels:
      env: prod
`)

	var j JobConfig
	err := yaml.Unmarshal(yamlData, &j)
	if err != nil {
		t.Fatalf("Failed to unmarshal JobConfig: %v", err)
	}

	if j.Name != "test_job" {
		t.Errorf("Expected Name 'test_job', got %q", j.Name)
	}
	if len(j.CollectorRefs) != 1 || j.CollectorRefs[0] != "collector1" {
		t.Errorf("Expected collectors [collector1], got %v", j.CollectorRefs)
	}
	if len(j.StaticConfigs) != 1 {
		t.Fatalf("Expected 1 static config, got %d", len(j.StaticConfigs))
	}
	if len(j.StaticConfigs[0].Targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(j.StaticConfigs[0].Targets))
	}
	if j.StaticConfigs[0].Labels["env"] != "prod" {
		t.Errorf("Expected label env=prod, got %q", j.StaticConfigs[0].Labels["env"])
	}
}

func TestMetricConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid gauge metric",
			yaml: `
metric_name: test_metric
type: gauge
help: Test help
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: false,
		},
		{
			name: "valid counter metric",
			yaml: `
metric_name: test_metric
type: counter
help: Test help
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: false,
		},
		{
			name: "missing name",
			yaml: `
type: gauge
help: Test help
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "missing type",
			yaml: `
metric_name: test_metric
help: Test help
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "missing help",
			yaml: `
metric_name: test_metric
type: gauge
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "missing values",
			yaml: `
metric_name: test_metric
type: gauge
help: Test help
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "unsupported type",
			yaml: `
metric_name: test_metric
type: histogram
help: Test help
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "multiple values without value_label",
			yaml: `
metric_name: test_metric
type: gauge
help: Test help
values: ["val1", "val2"]
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "multiple values with value_label",
			yaml: `
metric_name: test_metric
type: gauge
help: Test help
values: ["val1", "val2"]
value_label: metric_name
query: "SELECT 1"
`,
			wantErr: false,
		},
		{
			name: "duplicate key labels",
			yaml: `
metric_name: test_metric
type: gauge
help: Test help
key_labels: ["label1", "label1"]
values: ["val1"]
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "both query and query_ref",
			yaml: `
metric_name: test_metric
type: gauge
help: Test help
values: ["val1"]
query: "SELECT 1"
query_ref: some_query
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m MetricConfig
			err := yaml.Unmarshal([]byte(tt.yaml), &m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestQueryConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid query",
			yaml: `
query_name: test_query
query: "SELECT 1"
`,
			wantErr: false,
		},
		{
			name: "missing name",
			yaml: `
query: "SELECT 1"
`,
			wantErr: true,
		},
		{
			name: "missing query",
			yaml: `
query_name: test_query
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q QueryConfig
			err := yaml.Unmarshal([]byte(tt.yaml), &q)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestSecretMarshal(t *testing.T) {
	s := Secret("mysecret")
	result, err := s.MarshalYAML()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "<secret>" {
		t.Errorf("Expected '<secret>', got %q", result)
	}

	empty := Secret("")
	result, err = empty.MarshalYAML()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil for empty secret, got %v", result)
	}
}

func TestCollectorConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid collector with literal query",
			yaml: `
collector_name: test_collector
metrics:
  - metric_name: test_metric
    type: gauge
    help: Test help
    values: ["val1"]
    query: "SELECT 1"
`,
			wantErr: false,
		},
		{
			name: "no metrics",
			yaml: `
collector_name: test_collector
`,
			wantErr: true,
		},
		{
			name: "unresolved query_ref",
			yaml: `
collector_name: test_collector
metrics:
  - metric_name: test_metric
    type: gauge
    help: Test help
    values: ["val1"]
    query_ref: nonexistent_query
`,
			wantErr: true,
		},
		{
			name: "valid collector with named query",
			yaml: `
collector_name: test_collector
queries:
  - query_name: my_query
    query: "SELECT 1 AS val1"
metrics:
  - metric_name: test_metric
    type: gauge
    help: Test help
    values: ["val1"]
    query_ref: my_query
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c CollectorConfig
			err := yaml.Unmarshal([]byte(tt.yaml), &c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestJobConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid job",
			yaml: `
job_name: test_job
collectors: ["c1"]
static_configs:
  - targets:
      t1: "sqlserver://user:pass@localhost:1433"
`,
			wantErr: false,
		},
		{
			name: "missing name",
			yaml: `
collectors: ["c1"]
static_configs:
  - targets:
      t1: "sqlserver://user:pass@localhost:1433"
`,
			wantErr: true,
		},
		{
			name: "no static configs",
			yaml: `
job_name: test_job
collectors: ["c1"]
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var j JobConfig
			err := yaml.Unmarshal([]byte(tt.yaml), &j)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	configContent := `
global:
  scrape_timeout: 10s
  scrape_timeout_offset: 500ms
  min_interval: 0s
  max_connections: 3
  max_idle_connections: 3
target:
  data_source_name: "sqlserver://user:pass@localhost:1433"
  collectors: [test_collector]
collectors:
  - collector_name: test_collector
    metrics:
      - metric_name: test_metric
        type: gauge
        help: Test help
        values: ["val1"]
        query: "SELECT 1 AS val1"
`
	tmpfile, err := ioutil.TempFile("", "config.*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg == nil {
		t.Fatalf("Expected config, got nil")
	}
	if cfg.Target == nil {
		t.Fatalf("Expected target, got nil")
	}
	if cfg.Target.DSN != "sqlserver://user:pass@localhost:1433" {
		t.Errorf("Expected DSN 'sqlserver://user:pass@localhost:1433', got %q", cfg.Target.DSN)
	}
	if len(cfg.Collectors) != 1 {
		t.Errorf("Expected 1 collector, got %d", len(cfg.Collectors))
	}
	if cfg.Collectors[0].Name != "test_collector" {
		t.Errorf("Expected collector name 'test_collector', got %q", cfg.Collectors[0].Name)
	}
}

func TestLoadWithJobs(t *testing.T) {
	configContent := `
global:
  scrape_timeout: 10s
  scrape_timeout_offset: 500ms
  min_interval: 0s
  max_connections: 3
  max_idle_connections: 3
jobs:
  - job_name: test_job
    collectors: [test_collector]
    static_configs:
      - targets:
          target1: "sqlserver://user:pass@localhost:1433"
collectors:
  - collector_name: test_collector
    metrics:
      - metric_name: test_metric
        type: gauge
        help: Test help
        values: ["val1"]
        query: "SELECT 1 AS val1"
`
	tmpfile, err := ioutil.TempFile("", "config.*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cfg.Jobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Name != "test_job" {
		t.Errorf("Expected job name 'test_job', got %q", cfg.Jobs[0].Name)
	}
	if len(cfg.Jobs[0].Collectors()) != 1 {
		t.Errorf("Expected 1 resolved collector, got %d", len(cfg.Jobs[0].Collectors()))
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestConfigYAML(t *testing.T) {
	configContent := `
global:
  scrape_timeout: 10s
  scrape_timeout_offset: 500ms
  min_interval: 0s
  max_connections: 3
  max_idle_connections: 3
target:
  data_source_name: "sqlserver://user:pass@localhost:1433"
  collectors: [test_collector]
collectors:
  - collector_name: test_collector
    metrics:
      - metric_name: test_metric
        type: gauge
        help: Test help
        values: ["val1"]
        query: "SELECT 1 AS val1"
`
	tmpfile, err := ioutil.TempFile("", "config.*.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	yamlBytes, err := cfg.YAML()
	if err != nil {
		t.Fatalf("Expected no error from YAML(), got %v", err)
	}
	if len(yamlBytes) == 0 {
		t.Error("Expected non-empty YAML output")
	}
}

func TestTargetCollectors(t *testing.T) {
	target := &TargetConfig{
		DSN:        "sqlserver://user:pass@localhost:1433",
		collectors: []*CollectorConfig{{Name: "c1"}, {Name: "c2"}},
	}
	colls := target.Collectors()
	if len(colls) != 2 {
		t.Errorf("Expected 2 collectors, got %d", len(colls))
	}
}

func TestJobCollectors(t *testing.T) {
	j := &JobConfig{
		Name:       "test_job",
		collectors: []*CollectorConfig{{Name: "c1"}},
	}
	colls := j.Collectors()
	if len(colls) != 1 {
		t.Errorf("Expected 1 collector, got %d", len(colls))
	}
}

func TestGlobalConfigInvalidOffset(t *testing.T) {
	yamlData := []byte(`
scrape_timeout_offset: 0s
`)
	var g GlobalConfig
	err := yaml.Unmarshal(yamlData, &g)
	if err == nil {
		t.Error("Expected error for zero scrape_timeout_offset, got nil")
	}
}
