package test

import (
	"kis-folw/metrics"
	"testing"
)

func TestPrometheusServer(t *testing.T) {

	err := metrics.RunMetricsService("0.0.0.0:20004")
	if err != nil {
		panic(err)
	}
}
