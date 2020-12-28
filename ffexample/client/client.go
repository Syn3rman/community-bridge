package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	otelhttptrace "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

var logger = log.New(os.Stderr, "fluent-test", log.Ldate|log.Ltime|log.Llongfile)

func main() {
	tr := otel.Tracer("ffexample/client")
	ctx := baggage.ContextWithValues(context.Background(),
		label.String("foo", "bar"))
	ctx, span := tr.Start(ctx, "fib")
	defer span.End()
	test1(ctx, tr)
}

func test1(ctx context.Context, tr trace.Tracer) {

	url := "http://localhost:5050/fib"
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	ctx, childSpan := tr.Start(ctx, "inside test1")
	defer childSpan.End()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	check(err)
	_, req = otelhttptrace.W3C(ctx, req)
	fmt.Println("Sending request")
	res, err := client.Do(req)
	check(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	check(err)
	fmt.Printf("Response received (HTTP status code %d): %s\n", res.StatusCode, body)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
