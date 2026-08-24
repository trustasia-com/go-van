package grpcx

import (
	"testing"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/telemetry"
)

type testTelemetryRuntime struct {
	signals telemetry.Signal
}

func (runtime testTelemetryRuntime) Enabled(signal telemetry.Signal) bool {
	return runtime.signals.Enabled(signal)
}

func TestTelemetryServerOptionsFollowRuntimeSignals(t *testing.T) {
	var disabled *telemetry.Runtime
	tests := []struct {
		name    string
		runtime server.TelemetryRuntime
		want    int
	}{
		{name: "omitted", runtime: nil, want: 0},
		{name: "disabled", runtime: disabled, want: 0},
		{name: "tracer", runtime: testTelemetryRuntime{signals: telemetry.SignalTracer}, want: 1},
		{name: "meter", runtime: testTelemetryRuntime{signals: telemetry.SignalMeter}, want: 2},
		{name: "tracer and meter", runtime: testTelemetryRuntime{signals: telemetry.SignalTracer | telemetry.SignalMeter}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(telemetryServerOptions(tt.runtime)); got != tt.want {
				t.Fatalf("telemetryServerOptions() length = %d, want %d", got, tt.want)
			}
		})
	}
}
