package modem

import (
	"strings"

	modemconfig "github.com/iskrapw/modem/config"
	"github.com/iskrapw/utils/misc"
)

const _ModemInstantiationError = "unable to instantiate modem"
const _UnknownModulationError = "unknown modulation"
const _UnknownModemError = "unknown modem"

type NamedModem struct {
	modem  Modem
	config modemconfig.Modulation
}

type SwitchableModem struct {
	selectedModem   Modem
	availableModems []NamedModem
}

func (s *SwitchableModem) Initialize(config modemconfig.Config) error {
	s.availableModems = make([]NamedModem, len(config.Modulations))

	for i, modulationConfig := range config.Modulations {
		modem, err := instantiateModem(config, modulationConfig)
		if err != nil {
			return misc.WrapError(_ModemInstantiationError, err)
		}
		s.availableModems[i] = NamedModem{
			modem:  modem,
			config: modulationConfig,
		}
	}

	return s.Switch(config.InitialModulation)
}

func (s *SwitchableModem) Switch(modulationName string) error {
	for _, namedModem := range s.availableModems {
		if namedModem.config.Name == modulationName {
			s.selectedModem = namedModem.modem
			return nil
		}
	}

	return misc.NewError(_UnknownModulationError, modulationName)
}

func (s *SwitchableModem) Modulate(message []byte) []float64 {
	return s.selectedModem.Modulate(message)
}

func (s *SwitchableModem) Demodulate(samples []float64) []byte {
	return s.selectedModem.Demodulate(samples)
}

func instantiateModem(config modemconfig.Config, modulationConfig modemconfig.Modulation) (Modem, error) {
	modemName := strings.ToUpper(modulationConfig.Modem)
	inputSampleRate := config.Connections.Sound.InputSampleRate
	outputSampleRate := config.Connections.Sound.OutputSampleRate

	switch modemName {
	case "RTTY":
		return instantiateRTTYModem(config, modulationConfig)
	case "BELL103":
		bellModem := &Bell103Modem{}
		bellModem.Initialize(inputSampleRate, outputSampleRate, modulationConfig.Center)
		return bellModem, nil
	case "BELL202":
		bellModem := &Bell202Modem{}
		bellModem.Initialize(inputSampleRate, outputSampleRate, modulationConfig.Center)
		return bellModem, nil
	}

	return &RTTYModem{}, misc.NewError(_UnknownModemError, modemName)
}

func instantiateRTTYModem(config modemconfig.Config, modulationConfig modemconfig.Modulation) (*RTTYModem, error) {
	rttyModem := &RTTYModem{}
	speed := modulationConfig.Parameters["speed"].(float64)
	shift := modulationConfig.Parameters["shift"].(float64)
	uartParameters := modulationConfig.Parameters["uart"].(map[string]any)
	dataBits := int(uartParameters["dataBits"].(float64))
	stopBits := uartParameters["stopBits"].(float64)
	rttyModem.Initialize(speed, config.Connections.Sound.InputSampleRate, config.Connections.Sound.OutputSampleRate, modulationConfig.Center, shift, dataBits, stopBits)
	return rttyModem, nil
}
