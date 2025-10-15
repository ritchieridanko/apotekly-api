package tracer

import (
	"context"
	"fmt"
	"log"

	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Tracer struct {
	Cleanup func()
}

func NewProvider(cfg *configs.Config) (*Tracer, error) {
	ctx := context.Background()

	exp, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(cfg.Tracer.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.App.Name),
			),
		),
	)

	otel.SetTracerProvider(tp)

	log.Println("✅ initialized tracer")
	return &Tracer{func() { _ = tp.Shutdown(ctx) }}, nil
}
