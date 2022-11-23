package app

import (
	"log"

	"github.com/iskrapw/utils/misc"
)

func (app *ModemApplication) Close() error {
	err := app.soundClient.Disconnect()
	if err != nil {
		log.Println(err)
	}

	app.dataServer.Stop()

	log.Println("Closing threads")
	misc.BlockForSeconds(1)

	log.Println("Closing")
	return nil
}
