package metrics

import (
	"testing"
)

func TestGetMetricsInfo(t *testing.T) {
	// Create a simple test registry to check the slice approach
	registry := &Registry{}

	// Manually add some test metric info to verify the slice works
	registry.addMetricInfo("test_metric", "Test help text", []string{"label1", "label2"})
	registry.addMetricInfo("another_metric", "Another help text", []string{"label3"})

	metricsInfo := registry.GetMetricsInfo()

	if len(metricsInfo) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(metricsInfo))
	}

	if metricsInfo[0].Name != "test_metric" {
		t.Errorf("Expected first metric name 'test_metric', got '%s'", metricsInfo[0].Name)
	}

	if metricsInfo[0].Help != "Test help text" {
		t.Errorf("Expected first metric help 'Test help text', got '%s'", metricsInfo[0].Help)
	}
}
