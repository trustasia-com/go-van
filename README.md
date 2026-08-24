# go-van
Go framework for microservices.

See [Examples](./examples) for more usage.

## Telemetry

Starting OpenTelemetry providers is an application lifecycle concern. Enable it explicitly, then
select the matching server middleware independently:

```go
runtime, err := telemetry.Start(ctx,
	telemetry.WithName("example-api"),
	telemetry.WithEndpoint("localhost:4317"),
	telemetry.WithSignals(telemetry.SignalTracer),
	telemetry.WithInsecure(),
)
if err != nil {
	return err
}
defer runtime.Shutdown(shutdownCtx)

server := httpx.NewServer(
	server.WithHandler(handler),
	server.WithTelemetry(runtime),
)
```

`httpx` and `grpcx` automatically install the middleware and interceptors matching the Runtime's
Signals. When Telemetry is disabled, do not call `telemetry.Start`; pass a nil Runtime or omit
`server.WithTelemetry`. The server then installs no telemetry processing chain and opens no OTLP
connection.
