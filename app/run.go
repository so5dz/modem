package app

import "log"

func (app *ModemApplication) Run() error {
	log.Println("Connecting to sound server")
	err := app.soundClient.Connect()
	if err != nil {
		return err
	}

	log.Println("Starting data server")
	err = app.dataServer.Start()
	if err != nil {
		return err
	}

	log.Println("Starting extra server")
	err = app.extraServer.Start()
	if err != nil {
		return err
	}

	log.Println("MARDES-modem started")
	return nil
}
