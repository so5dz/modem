package main

import (
	"github.com/so5dz/modem/app"
	modemconfig "github.com/so5dz/modem/config"
	"github.com/so5dz/utils/config"
	"github.com/so5dz/utils/misc"
)

const _ConfigLoadError = "unable to read/parse configuration"
const _AppInitializationError = "unable to initialize application"
const _AppStartError = "unable to start application"

func main() {
	misc.WrapMain(mainWithError)()
}

func mainWithError() error {
	cfg, err := config.LoadConfigFromArgs[modemconfig.Config]()
	if err != nil {
		return misc.WrapError(_ConfigLoadError, err)
	}

	var app app.ModemApplication

	err = app.Initialize(cfg)
	if err != nil {
		return misc.WrapError(_AppInitializationError, err)
	}

	err = app.Run()
	if err != nil {
		return misc.WrapError(_AppStartError, err)
	}

	misc.BlockUntilInterrupted()

	return app.Close()
}
