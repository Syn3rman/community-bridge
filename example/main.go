package main

import (
	"context"
	"log"

	ff "github.com/Syn3rman/fluentforward"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tp *sdktrace.TracerProvider

// initTracer creates and registers trace provider instance.
func initTracer(url string) {
	var err error
	exp, err := ff.NewRawExporter(
		url,
		"fluentforward",
		30,
	)
	if err != nil {
		log.Panicf("failed to init exporter: %v\n", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithConfig(
			sdktrace.Config{
				DefaultSampler: sdktrace.AlwaysSample(),
				Resource:       resource.NewWithAttributes(label.String("service.name", "fluentbitexample")),
			},
		),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)

}

func main() {
	// initialize trace provider.
	initTracer("localhost:24224")

	// Create a named tracer with package path as its name.
	tracer := tp.Tracer("example/fluentforward/main")
	ctx := context.Background()
	defer func() { _ = tp.Shutdown(ctx) }()

	var span trace.Span
	ctx, span = tracer.Start(ctx, "test operation")
	defer span.End()
	span.AddEvent("Can add info here")

	if err := subOperation(ctx); err != nil {
		panic(err)
	}
}

func subOperation(ctx context.Context) error {
	tr := tp.Tracer("example/fluentforward/suboperation")

	var span trace.Span
	_, span = tr.Start(ctx, "Sub operation")
	defer span.End()
	span.AddEvent("This is a sub span")
	return nil

}
