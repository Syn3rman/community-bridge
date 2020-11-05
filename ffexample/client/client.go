package main

import (
	"fmt"
	"net/http"
	"context"
	"io/ioutil"
	"os"
	"log"

	"go.opentelemetry.io/otel/api/trace"
		"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otelhttptrace "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	ff "fluentforward"
)

var logger = log.New(os.Stderr, "fluent-test", log.Ldate|log.Ltime|log.Llongfile)

func initTracer(url string){
	err := ff.InstallNewPipeline(
		url,
		"ff client",
		ff.WithLogger(logger),
		ff.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)

	if err!=nil{
		log.Fatal(err)
	}
}

func main(){
	initTracer("http://localhost:5050/")
	tr := global.Tracer("ffexample/client")
	ctx := otel.ContextWithBaggageValues(context.Background(),
		label.String("n", "12"))
	ctx, span := tr.Start(ctx, "fib")
	defer span.End()
	test1(ctx, tr) 
}

func test1(ctx context.Context, tr trace.Tracer){
	
	url := "http://localhost:5050/fib"
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	ctx, childSpan := tr.Start(ctx, "inside test1")
	defer childSpan.End()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	check(err)
	_, req = otelhttptrace.W3C(ctx, req)
	fmt.Println("Sending request: ")
	res, err := client.Do(req)
	check(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	check(err)
	fmt.Printf("Response received (HTTP status code %d): %s\n", res.StatusCode, body)
}

func check(err error) {
	if err!=nil{
		log.Fatal(err)
	}
}
