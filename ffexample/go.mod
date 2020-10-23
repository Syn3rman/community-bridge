module example.com/ffexample

go 1.14

replace example.com/fluentforward => ./fluentforward

require (
	example.com/fluentforward v0.0.0-00010101000000-000000000000 // indirect
	github.com/ugorji/go v1.1.13 // indirect
	github.com/vmihailenco/msgpack/v5 v5.0.0-beta.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.13.0 // indirect
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
)
