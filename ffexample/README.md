### Fluentbit exporter

To set up the application, you need to have fluent-bit installed on your machine.

Steps to run the example:

1. Clone the repo using the command:

```
git clone https://github.com/Syn3rman/community-bridge.git
```

2. cd into the ffexample directory:

```
cd community-bridge/ffexample
```

3. Start the server:

```
go run server/server.go
```

4. Start the fluentbit instance:

```
fluent-bit -c fluent.conf
```
5. Run the client

```
go run client/client.go
```
