module example.com/ffexample

go 1.14

replace fluentforward => ./fluentforward

require (
	fluentforward v0.0.0-00010101000000-000000000000
	github.com/vmihailenco/msgpack/v5 v5.0.0-beta.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.13.0
	go.opentelemetry.io/otel v0.13.0
)
