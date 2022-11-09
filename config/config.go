package config

type Config struct {
	DataPort          int          `json:"dataPort"`
	ExtraPort         int          `json:"extraPort"`
	Connections       Connections  `json:"connections"`
	InitialModulation string       `json:"initialModulation"`
	Modulations       []Modulation `json:"modulations"`
}

type Connections struct {
	Sound   SoundConnection   `json:"sound"`
	Control ControlConnection `json:"control"`
}

type SoundConnection struct {
	Host             string  `json:"host"`
	Port             int     `json:"port"`
	InputSampleRate  float64 `json:"inputSampleRate"`
	OutputSampleRate float64 `json:"outputSampleRate"`
}

type ControlConnection struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

type Modulation struct {
	Name       string     `json:"name"`
	Modem      string     `json:"modem"`
	Center     float64    `json:"center"`
	Parameters Parameters `json:"parameters"`
}

type Parameters map[string]any
