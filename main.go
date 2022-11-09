package main

import (
	"log"
	"math/rand"

	modemconfig "github.com/iskrapw/modem/config"
	"github.com/iskrapw/modem/modem"
	"github.com/iskrapw/network/tcp"
	"github.com/iskrapw/utils/config"
	"github.com/iskrapw/utils/convert"
	"github.com/iskrapw/utils/misc"
)

func main() {
	misc.WrapMain(mainWithError)()
}

func mainWithError() error {
	cfg, err := config.LoadConfigFromArgs[modemconfig.Config]()
	if err != nil {
		return err
	}

	log.Println("Initializing modems")
	var switchableModem modem.SwitchableModem
	err = switchableModem.Initialize(cfg)
	if err != nil {
		return err
	}

	log.Println("Initializing data server")
	dataServer := tcp.NewServer(cfg.DataPort, tcp.TCPConnectionMode_Message)

	log.Println("Initializing receive sample buffer")
	var receivedSampleBuffer convert.ByteFloatBuffer
	receivedSampleBuffer.Initialize()

	log.Println("Initializing sound client")
	soundClient := tcp.NewClient(cfg.Connections.Sound.Host, cfg.Connections.Sound.Port, tcp.TCPConnectionMode_Stream)
	soundClient.OnReceive(func(b []byte) {
		receivedSampleBuffer.Put(b)
		samples := receivedSampleBuffer.GetAll()
		bytes := switchableModem.Demodulate(samples)
		if len(bytes) > 0 {
			dataServer.Broadcast(bytes)
		}
	})

	dataServer.OnReceive(func(client tcp.Remote, bytes []byte) {
		log.Println("Received", len(bytes), "bytes from client", client.Address())
		samples := switchableModem.Modulate(bytes)
		sampleBytes := convert.FloatsToBytes(samples)
		soundClient.Send(sampleBytes)
	})

	log.Println("Connecting to sound server")
	err = soundClient.Connect()
	if err != nil {
		return err
	}

	log.Println("Starting data server")
	err = dataServer.Start()
	if err != nil {
		return err
	}

	log.Println("MARDES-modem started, interrupt to close")

	misc.BlockUntilInterrupted()

	err = soundClient.Disconnect()
	if err != nil {
		log.Println(err)
	}

	dataServer.Stop()

	log.Println("Closing threads")
	misc.BlockForSeconds(1)

	log.Println("Closing")
	return nil
}

// todo remove
func randString(n int) string {
	var alphabet = []rune("qwertyuiopasdfghjklzxcvbnmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		v := rand.Intn(len(alphabet))
		s[i] = alphabet[v]
	}

	return string(s)
}
