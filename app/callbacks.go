package app

import (
	"log"

	"github.com/iskrapw/network/tcp"
	"github.com/iskrapw/utils/convert"
)

func (app *ModemApplication) onSoundReceived(sampleBytes []byte) {
	app.receivedSampleBuffer.Put(sampleBytes)
	samples := app.receivedSampleBuffer.GetAll()
	bytes := app.selectedModem.Demodulate(samples)
	if len(bytes) > 0 {
		app.dataServer.Broadcast(bytes)
	}
}

func (app *ModemApplication) onDataReceived(remote tcp.Remote, bytes []byte) {
	log.Println("[Data]  Received", len(bytes), "bytes from client", remote.Address())
	samples := app.selectedModem.Modulate(bytes)
	sampleBytes := convert.FloatsToBytes(samples)
	app.soundClient.Send(sampleBytes)
}

func (app *ModemApplication) onExtraCommand(remote tcp.Remote, bytes []byte) {
	log.Println("[Extra] Received", len(bytes), "bytes from client", remote.Address())
}
