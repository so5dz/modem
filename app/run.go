package app

import (
	"log"

	"github.com/so5dz/utils/misc"
)

const _SoundConnectionError = "unable to connect to sound server"
const _DataServerStartError = "unable to start data server"
const _ExtraServerStartError = "unable to start extra server"

func (app *ModemApplication) Run() error {
	log.Println("Connecting to sound server")
	err := app.soundClient.Connect()
	if err != nil {
		return misc.WrapError(_SoundConnectionError, err)
	}

	log.Println("Starting data server")
	err = app.dataServer.Start()
	if err != nil {
		return misc.WrapError(_DataServerStartError, err)
	}

	log.Println("Starting extra server")
	err = app.extraServer.Start()
	if err != nil {
		return misc.WrapError(_ExtraServerStartError, err)
	}

	log.Println("MARDES-modem started")
	return nil
}
