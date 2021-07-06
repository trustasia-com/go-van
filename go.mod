module github.com/trustasia-com/go-van

go 1.16

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/google/uuid v1.2.0
	github.com/justinas/alice v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/zouyx/agollo/v4 v4.0.7
	go.etcd.io/etcd/client/v3 v3.5.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.21.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC1
	go.opentelemetry.io/otel/metric v0.21.0
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/sdk/metric v0.21.0
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto v0.0.0-20210629200056-84d6f6074151
	google.golang.org/grpc v1.39.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
