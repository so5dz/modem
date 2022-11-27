package app

import (
	"log"

	"github.com/so5dz/modem/config"
)

func (app *ModemApplication) Initialize(cfg config.Config) error {
	app.initializeBuffers()
	app.initializeConnections(cfg)
	return app.initializeModems(cfg)
}

func (app *ModemApplication) initializeBuffers() {
	log.Println("Initializing receive sample buffer")
	app.receivedSampleBuffer.Initialize()
}

func (app *ModemApplication) initializeConnections(cfg config.Config) {
	log.Println("Initializing connections")
	app.initializeSoundClient(cfg)
	app.initializeDataServer(cfg)
	app.initializeExtraServer(cfg)
}

func (app *ModemApplication) initializeSoundClient(cfg config.Config) {
	log.Println("Initializing sound client")
	app.soundClient.Initialize(cfg.Connections.Sound.Host, cfg.Connections.Sound.Port)
	app.soundClient.OnReceive(app.onSoundReceived)
}

func (app *ModemApplication) initializeDataServer(cfg config.Config) {
	log.Println("Initializing data server")
	app.dataServer.Initialize(cfg.DataPort)
	app.dataServer.OnReceive(app.onDataReceived)
}

func (app *ModemApplication) initializeExtraServer(cfg config.Config) {
	log.Println("Initializing extra server")
	app.extraServer.Initialize(cfg.ExtraPort)
	app.extraServer.OnReceive(app.onExtraCommand)
}
