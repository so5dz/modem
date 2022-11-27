package app

import (
	"github.com/so5dz/modem/config"
	"github.com/so5dz/modem/modem"
	tcpc "github.com/so5dz/network/client/tcp"
	tcps "github.com/so5dz/network/server/tcp"
	"github.com/so5dz/utils/convert"
)

const _ModemInstantiationError = "unable to instantiate modem"
const _UnknownModulationError = "unknown modulation"
const _UnknownModemError = "unknown modem"

type namedModem struct {
	modem  modem.Modem
	config config.Modulation
}

type ModemApplication struct {
	selectedModem        modem.Modem
	availableModems      []namedModem
	soundClient          tcpc.StreamClient
	dataServer           tcps.StreamServer
	extraServer          tcps.StreamServer
	receivedSampleBuffer convert.ByteFloatBuffer
}
