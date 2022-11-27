package app

import (
	"log"

	"github.com/so5dz/network/server"
	"github.com/so5dz/utils/convert"
)

func (app *ModemApplication) onSoundReceived(sampleBytes []byte) {
	app.receivedSampleBuffer.Put(sampleBytes)
	samples := app.receivedSampleBuffer.GetAll()
	bytes := app.selectedModem.Demodulate(samples)
	if len(bytes) > 0 {
		app.dataServer.Broadcast(bytes)
	}
}

func (app *ModemApplication) onDataReceived(remote server.Remote, bytes []byte) {
	log.Println("[Data]  Received", len(bytes), "bytes from client", remote.Address())
	samples := app.selectedModem.Modulate(bytes)
	sampleBytes := convert.FloatsToBytes(samples)
	app.soundClient.Send(sampleBytes)
}

func (app *ModemApplication) onExtraCommand(remote server.Remote, bytes []byte) {
	log.Println("[Extra] Received", len(bytes), "bytes from client", remote.Address())
}
