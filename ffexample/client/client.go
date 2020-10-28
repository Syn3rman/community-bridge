package main

import (
	"fmt"
	"net/http"
	"context"
	"io/ioutil"
	"os"

		"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otelhttptrace "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	ff "fluentforward"
)

func initTracer(url string){
	err := ff.InstallNewPipeline(url, "ff client")
	if err!=nil{
		fmt.Println(err)
	}
}

func main(){
	initTracer("http://localhost:5050/")
	tr := global.Tracer("ffexample/client")
	url := "http://localhost:5050/fib"
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	ctx := otel.ContextWithBaggageValues(context.Background(),
		label.String("n", "12"))
	ctx, span := tr.Start(ctx, "fib")
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil{
		fmt.Println(err)
	}
	_, req = otelhttptrace.W3C(ctx, req)
	fmt.Println("Sending request: ")
	res, err := client.Do(req)
	if err!=nil{
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read HTTP response body: %v\n", err)
	}

	fmt.Printf("Response received (HTTP status code %d): %s\n", res.StatusCode, body)
}
