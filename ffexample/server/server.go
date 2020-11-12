package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ff "fluentforward"
	otelhttptrace "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

func initTracer(url string) {
	err := ff.InstallNewPipeline(url, "fluentforward")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	initTracer("localhost:24224")

	tr := global.Tracer("ffexample/server")

	fib := func(w http.ResponseWriter, req *http.Request) {
		n := req.FormValue("n")
		attrs, _, spanCtx := otelhttptrace.Extract(req.Context(), req)

		_, span := tr.Start(
			trace.ContextWithRemoteSpanContext(req.Context(), spanCtx),
			"hello",
			trace.WithAttributes(attrs...),
			trace.WithLinks(trace.Link{SpanContext: spanCtx, Attributes: attrs}),
		)
		ctx := context.Background()
		ctx = otel.ContextWithBaggageValues(ctx, label.String("foo2", "foo1"), label.String("bar1", "bar3"))
		span.AddEvent(ctx, "testEvent", label.String("New", "attr"))
		span.SetStatus(2, "setting span status")
		defer span.End()

		json.NewEncoder(w).Encode(n)
	}

	http.HandleFunc("/fib", fib)
	http.ListenAndServe(":5050", nil)
}
