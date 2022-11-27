package main

import (
	"github.com/so5dz/modem/app"
	modemconfig "github.com/so5dz/modem/config"
	"github.com/so5dz/utils/config"
	"github.com/so5dz/utils/misc"
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
