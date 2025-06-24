package main

import (
	"context"
	"log"
	"time"

	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func main() {
	ctx := context.Background()

	// 初始化telemetry，添加自定义变量
	shutdown, _ := telemetry.InitProvider(ctx,
		telemetry.WithName("custom-service"),
		telemetry.WithEndpoint("localhost:4317"),
		telemetry.WithFlag(telemetry.FlagTracer|telemetry.FlagMeter|telemetry.FlagLogger|telemetry.FlagInsecure),
		telemetry.WithAttributes(
			attribute.String("version", "1.0.0"),
			attribute.String("environment", "production"),
			attribute.String("region", "us-west-2"),
			attribute.String("team", "backend"),
			attribute.String("component", "user-service"),
			attribute.Int("instance_id", 12345),
			semconv.K8SNamespaceNameKey.String("default"),
		),
	)
	defer shutdown()

	// 使用tracer - 自定义属性会自动包含在所有span中
	tracer := otel.Tracer("custom-service")
	ctx, span := tracer.Start(ctx, "example-operation")
	defer span.End()

	// 为当前span添加额外的属性
	span.SetAttributes(
		attribute.String("user_id", "user123"),
		attribute.String("request_id", "req456"),
	)

	// 使用meter - 自定义属性会自动包含在所有metrics中
	meter := otel.Meter("custom-service")
	counter, err := meter.Int64Counter(
		"requests_total",
		metric.WithDescription("Total number of requests"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 记录metrics，可以添加额外的标签
	counter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", "GET"),
			attribute.String("endpoint", "/api/users"),
			attribute.Int("status_code", 200),
		),
	)

	// 模拟一些工作
	time.Sleep(100 * time.Millisecond)

	logx.Info("Telemetry example completed")
}
