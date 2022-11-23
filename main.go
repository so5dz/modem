package main

import (
	"github.com/iskrapw/modem/app"
	modemconfig "github.com/iskrapw/modem/config"
	"github.com/iskrapw/utils/config"
	"github.com/iskrapw/utils/misc"
)

func main() {
	misc.WrapMain(mainWithError)()
}

func mainWithError() error {
	cfg, err := config.LoadConfigFromArgs[modemconfig.Config]()
	if err != nil {
		return err
	}

	var app app.ModemApplication

	err = app.Initialize(cfg)
	if err != nil {
		return err
	}

	err = app.Run()
	if err != nil {
		return err
	}

	misc.BlockUntilInterrupted()

	return app.Close()
}
