package fluentforward

import (
	"fmt"
	"net"
	"time"
	"context"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/api/global"
)

type ffSpan struct{
	_msgpack struct{} `msgpack:",asArray"`
	Name string `msgpack:"name"`
	Ts int64 `msgpack:"ts"`
	Attrs map[label.Key]interface{} `msgpack:"attrs"`
}

type Exporter struct{
	url string
	serviceName string
	client *http.Client
}

func InstallNewPipeline(ffurl, serviceName string) error{
	tp,err := NewExportPipeline(ffurl, serviceName)
	if err!=nil{
		fmt.Println(err)
	}
	global.SetTracerProvider(tp)
	return nil
}

func NewExportPipeline(ffurl, serviceName string) (*sdktrace.TracerProvider, error){
	exp, err := NewExporter(ffurl, serviceName)
	if err!=nil{
		fmt.Println(err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
	return tp, nil
}

func NewExporter(ffurl, serviceName string) (*Exporter, error){

	return &Exporter{
		url: ffurl,
		client: http.DefaultClient,
		serviceName: serviceName,
	}, nil

}

func (e *Exporter) ExportSpans(ctx context.Context, sds []*export.SpanData) error{
	fmt.Println("Exporting spans to fluentd")
	spans := ffSpan{
		Name: "span",
		Ts: time.Now().UnixNano(),
	}
	testkv := make(map[label.Key]interface{})
	for _, span := range sds{
		for _, val := range span.Attributes{
			testkv[val.Key] = val.Value.AsInterface()
		}
	}

	spans.Attrs = testkv

	t, err := msgpack.Marshal(&spans)
	check(err)

	url := "localhost:24224"
	tcpAddr, err := net.ResolveTCPAddr("tcp", url)
	check(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	check(err)

	_, err = conn.Write(t)
	check(err)
	return nil
}

func (e *Exporter) Shutdown(ctx context.Context) error{
	return nil
}

func check(err error){
	if err!=nil{
		fmt.Println(err)
	}
}
