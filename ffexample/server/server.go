package main 

import (
	"fmt"
	"net/http"
	"encoding/json"
	
	otelhttptrace "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	// "go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/global"
	// ff "github.com/Syn3rman/opentelemetry-go-contrib/exporters/trace/fluentforward"
	ff "example.com/fluentforward"
)

func initTracer(url string){
	err := ff.InstallNewPipeline(url, "fluentforward")
	if err!=nil{
		fmt.Println(err)
	}
}

func main(){
	initTracer("http://localhost:24224/")

	tr := global.Tracer("ffexample/server")

	fib := func(w http.ResponseWriter, req *http.Request){
		n := req.FormValue("n")
		attrs, _, spanCtx := otelhttptrace.Extract(req.Context(), req)

		fmt.Printf("Attributes: %+v\n\nSpan Context: %+v\n\n\n", attrs, spanCtx)
		_, span := tr.Start(
						trace.ContextWithRemoteSpanContext(req.Context(), spanCtx),
						"hello",
						trace.WithAttributes(attrs...),
						trace.WithLinks(trace.Link{SpanContext: spanCtx, Attributes: attrs}),
		)
		defer span.End()

		// span.SetAttributes(label.String("test setting span attrs", "yes"))

		json.NewEncoder(w).Encode(n)
	}

	http.HandleFunc("/fib", fib)
	http.ListenAndServe(":5050", nil)	
}
