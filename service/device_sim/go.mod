module github.com/charmingruby/devicio/service/device_sim

go 1.23.2

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/charmingruby/devicio/lib v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	google.golang.org/protobuf v1.36.6
)

require github.com/streadway/amqp v1.1.0 // indirect

replace github.com/charmingruby/devicio/lib => ../../lib
