module github.com/deepzz0/go-van/examples

go 1.15

replace github.com/deepzz0/go-van => ../

require (
	github.com/deepzz0/go-van v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.20.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)
