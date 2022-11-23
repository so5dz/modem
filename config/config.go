package config

type Config struct {
	DataPort          int          `json:"dataPort"`
	ExtraPort         int          `json:"extraPort"`
	Connections       Connections  `json:"connections"`
	InitialModulation string       `json:"initialModulation"`
	Modulations       []Modulation `json:"modulations"`
}

type Connections struct {
	Sound SoundConnection `json:"sound"`
}

type SoundConnection struct {
	Host             string  `json:"host"`
	Port             int     `json:"port"`
	InputSampleRate  float64 `json:"inputSampleRate"`
	OutputSampleRate float64 `json:"outputSampleRate"`
}

type Modulation struct {
	Name       string     `json:"name"`
	Modem      string     `json:"modem"`
	Center     float64    `json:"center"`
	Parameters Parameters `json:"parameters"`
}

type Parameters map[string]any
