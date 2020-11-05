package fluentforward

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
	"github.com/vmihailenco/msgpack/v5"
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
	StatusMessage string            `msgpack:"statusMessage"`
	EndTime int64                   `msgpack:"endTime"`
	Attrs map[label.Key]interface{} `msgpack:"attrs"`
	DroppedAttributeCount int       `msgpack:"droppedAttributesCount"`
	DroppedMessageEventCount int    `msgpack:"droppedMessageEventCount"`
	DroppedLinkCount int            `msgpack:"droppedLinkCount"`
	StatusCode string               `msgpack:"statusCode"`
	Links []link                    `msgpack:"links"`
}

type spanContext struct{
	TraceId string  `msgpack:"traceId"`
	SpanId string   `msgpack:"spanId"`
	TraceFlags byte `msgpack:"TraceFlags"`
}

type link struct{
	SpanContext spanContext         `msgpack:"spanContext"`
	Attrs map[label.Key]interface{} `msgpack:"attrs"`
}

type event struct{
	Attrs map[label.Key]interface{} `msgpack:"attrs"`
	DroppedAttributeCount int       `msgpack:"droppedAttributesCount"`
}

type Exporter struct{
	url string
	serviceName string
	client *net.TCPConn
	o options
}

type Option func(*options)

type options struct{
	config *sdktrace.Config
	logger *log.Logger
}

func WithLogger(logger *log.Logger) Option{
	return func(o *options)	{
		o.logger = logger
	}
}

func WithSDK(config *sdktrace.Config) Option{
	return func(o *options){
		o.config = config
	}
}

func InstallNewPipeline(ffurl, serviceName string, opts ...Option) error{
	tp,err := NewExportPipeline(ffurl, serviceName, opts...)
	if err!=nil{
		return err
	}
	global.SetTracerProvider(tp)
	return nil
}

func NewExportPipeline(ffurl, serviceName string, opts ...Option) (*sdktrace.TracerProvider, error){
	exp, err := NewExporter(ffurl, serviceName, opts...)
	if err!=nil{
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
	return tp, nil
}

func NewExporter(ffurl, serviceName string, opts ...Option) (*Exporter, error){
	o := options{}
	for _, opt := range opts{
		opt(&o)
	}

	if ffurl == ""{
		return nil, errors.New("Fluent instance url cannot be empty")
	}

	return &Exporter{
		url: ffurl,
		serviceName: serviceName,
		o: o,
	}, nil

}

func (e *Exporter) ExportSpans(ctx context.Context, sds []*export.SpanData) error{

	ts := time.Now().Unix()

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
		spans.StatusMessage = span.StatusMessage
		spans.StatusCode = span.StatusCode.String()
		spans.StartTime = span.StartTime.UnixNano()
		spans.EndTime = span.EndTime.UnixNano()
		spans.DroppedAttributeCount = span.DroppedAttributeCount
		spans.DroppedLinkCount = span.DroppedLinkCount
		spans.DroppedMessageEventCount = span.DroppedMessageEventCount
		
		spans.Attrs = attributesToMap(span.Attributes)
	
		spans.Links = linksToSlice(span.Links)

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

func attributesToMap(attributes []label.KeyValue) map[label.Key]interface{}{
		attrs := make(map[label.Key]interface{})
		for _,v := range attributes{
			attrs[v.Key] = v.Value.AsInterface()
		}
		return attrs
}

func linksToSlice(links []apitrace.Link) []link{
	var l []link
	for _, v := range links{
		temp := link{
			SpanContext: spanContext{
				TraceId: v.SpanContext.TraceID.String(),
				SpanId: v.SpanContext.SpanID.String(),
				TraceFlags: v.SpanContext.TraceFlags,
			},
			Attrs: attributesToMap(v.Attributes),
		}
		l = append(l, temp)	
	}
	return l
}

func check(err error){
	if err!=nil{
		panic(err)
	}
}
