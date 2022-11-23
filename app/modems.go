package app

import (
	"log"
	"strings"

	"github.com/iskrapw/modem/config"
	"github.com/iskrapw/modem/modem"
	"github.com/iskrapw/utils/misc"
)

func (app *ModemApplication) SetModem(name string) error {
	for _, namedModem := range app.availableModems {
		if namedModem.config.Name == name {
			app.selectedModem = namedModem.modem
			return nil
		}
	}
	return misc.NewError(_UnknownModulationError, name)
}

func (app *ModemApplication) initializeModems(cfg config.Config) error {
	log.Println("Initializing modems")
	app.availableModems = make([]namedModem, len(cfg.Modulations))

	for i, modulationConfig := range cfg.Modulations {
		modem, err := instantiateModem(cfg, modulationConfig)
		if err != nil {
			return misc.WrapError(_ModemInstantiationError, err)
		}
		app.availableModems[i] = namedModem{
			modem:  modem,
			config: modulationConfig,
		}
	}

	return app.SetModem(cfg.InitialModulation)
}

func instantiateModem(cfg config.Config, modulationConfig config.Modulation) (modem.Modem, error) {
	modemName := strings.ToUpper(modulationConfig.Modem)
	inputSampleRate := cfg.Connections.Sound.InputSampleRate
	outputSampleRate := cfg.Connections.Sound.OutputSampleRate

	switch modemName {
	case "RTTY":
		return instantiateRTTYModem(cfg, modulationConfig)
	case "BELL103":
		bellModem := &modem.Bell103Modem{}
		bellModem.Initialize(inputSampleRate, outputSampleRate, modulationConfig.Center)
		return bellModem, nil
	case "BELL202":
		bellModem := &modem.Bell202Modem{}
		bellModem.Initialize(inputSampleRate, outputSampleRate, modulationConfig.Center)
		return bellModem, nil
	}

	return &modem.RTTYModem{}, misc.NewError(_UnknownModemError, modemName)
}

func instantiateRTTYModem(cfg config.Config, modulationConfig config.Modulation) (*modem.RTTYModem, error) {
	rttyModem := &modem.RTTYModem{}
	speed := modulationConfig.Parameters["speed"].(float64)
	shift := modulationConfig.Parameters["shift"].(float64)
	uartParameters := modulationConfig.Parameters["uart"].(map[string]any)
	dataBits := int(uartParameters["dataBits"].(float64))
	stopBits := uartParameters["stopBits"].(float64)
	rttyModem.Initialize(speed, cfg.Connections.Sound.InputSampleRate, cfg.Connections.Sound.OutputSampleRate, modulationConfig.Center, shift, dataBits, stopBits)
	return rttyModem, nil
}
