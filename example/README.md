### FluentBit exporter

This repo contains an example application that showcases using the FluentBit exporter to export span data to the opentelemetry collector, where it can be parsed enabling the `fluentforwardextreceiver` in the traces pipeline.

To run the example:

1.  Clone the repo

```sh
; git clone https://github.com/Syn3rman/community-bridge.git
```

2. Change directory

```sh
; cd community-bridge/example
```

3. Install dependencies

```go
; go get
```

4. Run example

```go
; go run main.go
```