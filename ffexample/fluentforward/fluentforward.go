package fluentforward

import (
	"fmt"
	"net"
	"time"
	"context"

	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/label"
	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/api/global"
)

type ffSpan struct{
	_msgpack struct{}            `msgpack:",asArray"`
	Tag string                   `msgpack:"tag"`
	Ts int64                     `msgpack:"ts"`
	SpanData spanData						 `msgpack:"spanData"`
}

type spanData struct{
	SpanContext spanContext         `msgpack:"spanContext"`
	ParentSpanId string             `msgpack:"parentSpanId"`
	SpanKind apitrace.SpanKind      `msgpack:"spanKind"`
	Name string                     `msgpack:"name"`
	StartTime int64                 `msgpack:"startTime"`
	EndTime int64                   `msgpack:"endTime"`
	Attrs map[label.Key]interface{} `msgpack:"attrs"`
}

type spanContext struct{
	TraceId string  `msgpack:"traceId"`
	SpanId string   `msgpack:"spanId"`
	TraceFlags byte `msgpack:"TraceFlags"`
}

type Exporter struct{
	url string
	serviceName string
	client *net.TCPConn
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
		serviceName: serviceName,
	}, nil

}

func (e *Exporter) ExportSpans(ctx context.Context, sds []*export.SpanData) error{
	ts := time.Now().Unix()
	fmt.Println("Timestamp: ", ts)

	tcpAddr, err := net.ResolveTCPAddr("tcp", e.url)
	check(err)
	client, err := net.DialTCP("tcp", nil, tcpAddr)
	e.client = client
	check(err)

	for _, span := range sds{
		ffspan := ffSpan{
			Tag: "span.test",
			Ts: ts,
		}
		spans := spanData{}
		spans.SpanContext = spanContext{
			TraceId: span.SpanContext.TraceID.String(),
			SpanId: span.SpanContext.SpanID.String(),
			TraceFlags: span.SpanContext.TraceFlags,
		}
		spans.ParentSpanId = span.ParentSpanID.String()
		spans.SpanKind = span.SpanKind
		spans.Name = span.Name
		spans.StartTime = span.StartTime.UnixNano()
		spans.EndTime = span.EndTime.UnixNano()
		
		testkv := make(map[label.Key]interface{})
		for _,v := range span.Attributes{
			testkv[v.Key] = v.Value.AsInterface()
		}
		spans.Attrs = testkv

		ffspan.SpanData = spans

		t, err := msgpack.Marshal(&ffspan)
		check(err)

		_, err = e.client.Write(t)
		check(err)
	}
	return nil
}

func (e *Exporter) Shutdown(ctx context.Context) error{
	e.client.Close()
	return nil
}

func check(err error){
	if err!=nil{
		panic(err)
	}
}
