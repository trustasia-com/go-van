package server

import (
	"testing"

	"github.com/trustasia-com/go-van/pkg/telemetry"
)

type testTelemetryRuntime struct {
	signals telemetry.Signal
}

func (runtime testTelemetryRuntime) Enabled(signal telemetry.Signal) bool {
	return runtime.signals.Enabled(signal)
}

func TestWithTelemetryStoresRuntimeAsSingleSignalSource(t *testing.T) {
	runtime := testTelemetryRuntime{signals: telemetry.SignalMeter | telemetry.SignalTracer}
	var options ServerOptions
	WithTelemetry(runtime)(&options)

	if options.Telemetry == nil {
		t.Fatal("WithTelemetry() did not store Runtime")
	}
	if !options.Telemetry.Enabled(telemetry.SignalMeter) ||
		!options.Telemetry.Enabled(telemetry.SignalTracer) {
		t.Fatal("ServerOptions Telemetry does not expose configured Runtime Signals")
	}
}
