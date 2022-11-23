package app

import (
	"github.com/iskrapw/modem/config"
	"github.com/iskrapw/modem/modem"
	"github.com/iskrapw/network/tcp"
	"github.com/iskrapw/utils/convert"
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
	soundClient          tcp.Client
	dataServer           tcp.Server
	extraServer          tcp.Server
	receivedSampleBuffer convert.ByteFloatBuffer
}
