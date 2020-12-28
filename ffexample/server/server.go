package main

import (
	"fmt"
	"io"
	"net/http"

	ff "github.com/Syn3rman/fluentforward"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func initTracer(url string) {
	err := ff.InstallNewPipeline(url, "fluentforward")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	initTracer("localhost:24224")

	tr := otel.Tracer("ffexample/server")

	http.HandleFunc("/fib", helloHandler)
	http.ListenAndServe(":5050", nil)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	tracer := trace.GlobalTracer()

	// Extracts the conventional HTTP span attributes,
	// distributed context tags, and a span context for
	// tracing this request.
	attrs, tags, spanCtx := httptrace.Extract(req.Context(), req)

	// Apply the distributed context tags to the request
	// context.
	req = req.WithContext(tag.WithMap(req.Context(), tag.NewMap(tag.MapUpdate{
		MultiKV: tags,
	})))

	// Start the server-side span, passing the remote
	// child span context explicitly.
	_, span := tracer.Start(
		req.Context(),
		"hello",
		trace.WithAttributes(attrs...),
		trace.ChildOf(spanCtx),
	)
	defer span.End()

	_, _ = io.WriteString(w, "Hello, world!\n")
}
