package main

import (
	"os"

	"github.com/charmingruby/devicio/lib/pkg/messaging/nats"
	"github.com/charmingruby/devicio/service/device/pkg/logger"
)

func main() {
	logger.New()

	_, err := nats.New(nil)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

}
