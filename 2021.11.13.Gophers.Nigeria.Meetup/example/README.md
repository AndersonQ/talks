 # Monitoring applications with OpenTelemetry: a practical example
 
## Requirements
 - Go 1.17+
 - Docker
 - Docker compose
 - Make

## Run the dependencies

```shell
make compose-dependencies 
```

## Run the applications

All the applications are instrumented and exporting the data to Jaeger, which is
started when running the dependencies.

```shell
make run-root
make run-http
make run-events
```

 - `cmd/root`

Is an HTTP server listening on port `4242`. It makes a HTTP request to `cmd/http`

- `cmd/http`

Is an HTTP server and a Kafka producer, the HTTP server runs on port `1618` by default.
There are two endpoints:
 - `/`: produces an event to the kafka topic `otel-example`
 - `/final`: logs and responds with the inbound request headers 

 - `cmd/events`:

Is a Kafka consumer which consumer from the topic `otel-example` and then calls `http://localhost:1618/final`

After all applications are running, call [http://localhost:4242](http://localhost:4242) a few times.
Check the logs, and head to [Jaeger UI](http://localhost:16686/search) to see the traces.

## Clean up

- To stop the containers:

```shell
make compose-down
```

- To clean up, erase the data and remove volumes:

```shell
make clean
```
