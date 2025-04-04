module github.com/charmingruby/devicio/service/device_sim

go 1.23.2

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/charmingruby/devicio/lib v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/charmingruby/devicio/lib => ../../lib
