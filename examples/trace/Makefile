.PHONY: proto

grpcsrv:
	@go run grpc-service/main.go

httpsrv:
	@go run http-service/main.go

httpint:
	@go run http-interface/main.go

http2http:
	@curl -v localhost:9000/http-to-http/1

http2grpc:
	@curl -v localhost:9000/http-to-grpc/1

http2grpc2http:
	@curl -v localhost:9000/http-to-grpc-to-http/1

http2grpc2httperr:
	@curl -v localhost:9000/http-to-grpc-to-http/2

# go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest
# mock generate data
traces:
	@telemetrygen traces --otlp-insecure --duration 5s
metrics:
	@telemetrygen metrics --duration 5s --otlp-insecure
logs:
	@telemetrygen logs --duration 5s --otlp-insecure

