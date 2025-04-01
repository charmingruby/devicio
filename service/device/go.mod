module github.com/charmingruby/devicio/service/device

go 1.23.2

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/charmingruby/devicio/lib v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nats.go v1.40.1 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/charmingruby/devicio/lib => ../../lib
