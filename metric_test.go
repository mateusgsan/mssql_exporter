package sql_exporter

import (
	"testing"

	"github.com/mateusgsan/mssql_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/golang/protobuf/proto"
)

func TestAutomaticMetricDesc(t *testing.T) {
	constLabels := []*dto.LabelPair{
		{Name: proto.String("job"), Value: proto.String("test_job")},
	}
	labels := []string{"label1", "label2"}

	desc := NewAutomaticMetricDesc("test_ctx", "test_metric", "test help", prometheus.GaugeValue, constLabels, labels...)

	if desc.Name() != "test_metric" {
		t.Errorf("Expected Name 'test_metric', got %q", desc.Name())
	}
	if desc.Help() != "test help" {
		t.Errorf("Expected Help 'test help', got %q", desc.Help())
	}
	if desc.ValueType() != prometheus.GaugeValue {
		t.Errorf("Expected ValueType GaugeValue, got %v", desc.ValueType())
	}
	if desc.LogContext() != "test_ctx" {
		t.Errorf("Expected LogContext 'test_ctx', got %q", desc.LogContext())
	}
	if len(desc.Labels()) != 2 || desc.Labels()[0] != "label1" || desc.Labels()[1] != "label2" {
		t.Errorf("Expected Labels [label1, label2], got %v", desc.Labels())
	}
	if len(desc.ConstLabels()) != 1 || desc.ConstLabels()[0].GetName() != "job" {
		t.Errorf("Expected ConstLabels [job], got %v", desc.ConstLabels())
	}
}

func TestNewMetric(t *testing.T) {
	desc := NewAutomaticMetricDesc("test_ctx", "test_metric", "test help", prometheus.GaugeValue, nil, "label1")
	
	metric := NewMetric(desc, 42.0, "val1")
	
	if metric.Desc() != desc {
		t.Errorf("Expected Desc to match")
	}

	var out dto.Metric
	err := metric.Write(&out)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if out.Gauge == nil || out.Gauge.GetValue() != 42.0 {
		t.Errorf("Expected Gauge value 42.0, got %v", out.Gauge)
	}
	if len(out.Label) != 1 || out.Label[0].GetName() != "label1" || out.Label[0].GetValue() != "val1" {
		t.Errorf("Expected Label [label1=val1], got %v", out.Label)
	}
}

func TestNewInvalidMetric(t *testing.T) {
	err := NewInvalidMetric(nil)
	if err.Desc() != nil {
		t.Errorf("Expected nil Desc, got %v", err.Desc())
	}
	if err.Write(nil) != nil {
		t.Errorf("Expected nil error, got %v", err.Write(nil))
	}
}

func TestMetricFamily(t *testing.T) {
	mc := &config.MetricConfig{
		Name: "test_metric",
		Help: "test help",
		Values: []string{"val_col"},
		KeyLabels: []string{"key_col"},
	}
	// We need to set the valueType, but it's unexported in config.MetricConfig.
	// However, we can test the basic initialization.
	
	constLabels := []*dto.LabelPair{
		{Name: proto.String("job"), Value: proto.String("test_job")},
	}

	mf, err := NewMetricFamily("test_ctx", mc, constLabels)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mf.Name() != "test_metric" {
		t.Errorf("Expected Name 'test_metric', got %q", mf.Name())
	}
	if mf.Help() != "test help" {
		t.Errorf("Expected Help 'test help', got %q", mf.Help())
	}
	if len(mf.Labels()) != 1 || mf.Labels()[0] != "key_col" {
		t.Errorf("Expected Labels [key_col], got %v", mf.Labels())
	}
}
