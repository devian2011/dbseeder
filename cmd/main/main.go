package main

import (
	"context"
	"dbseeder/internal"
	"flag"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
)

func main() {
	command := flag.String("command", "help", "seed database data")
	filePath := flag.String("schema", "./config/db.conf.yml", "Path to schema config")
	flag.Parse()

	ctx := context.Background()
	signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT)

	app, err := internal.NewApplication(ctx, *filePath)
	if err != nil {
		logrus.Errorln("Error on init application. Error: ", err.Error())
		return
	}
	if runErr := app.Run(*command); runErr != nil {
		logrus.Errorln("Error on run application. Error: ", runErr.Error())
	}
}
