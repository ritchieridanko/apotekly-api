package ot

import (
	"context"
	"log"

	"github.com/ritchieridanko/apotekly-api/user/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Tracer struct {
	Cleanup func()
}

func Initialize() (tracer *Tracer, err error) {
	ctx := context.Background()

	exp, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(config.TracerGetEndpoint()),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(config.AppGetName()),
			),
		),
	)

	otel.SetTracerProvider(tp)

	log.Println("SUCCESS -> initialized tracer")
	return &Tracer{func() { _ = tp.Shutdown(ctx) }}, nil
}
