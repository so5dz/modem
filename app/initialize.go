package app

import (
	"log"

	"github.com/iskrapw/modem/config"
	"github.com/iskrapw/network/tcp"
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
	app.soundClient = tcp.NewClient(cfg.Connections.Sound.Host, cfg.Connections.Sound.Port, tcp.TCPConnectionMode_Stream)
	app.soundClient.OnReceive(app.onSoundReceived)
}

func (app *ModemApplication) initializeDataServer(cfg config.Config) {
	log.Println("Initializing data server")
	app.dataServer = tcp.NewServer(cfg.DataPort, tcp.TCPConnectionMode_Stream)
	app.dataServer.OnReceive(app.onDataReceived)
}

func (app *ModemApplication) initializeExtraServer(cfg config.Config) {
	log.Println("Initializing extra server")
	app.extraServer = tcp.NewServer(cfg.ExtraPort, tcp.TCPConnectionMode_Message)
	app.extraServer.OnReceive(app.onExtraCommand)
}
