package main

import "github.com/Anveena/opensslreproducer/quic"

func main() {
	config := &quic.Config{
		Port:                   45678,
		ReadChanSizeForClient:  16,
		WriteChanSizeForClient: 16,
		ReadBufferSize:         -1,
		WriteBufferSize:        -1,
	}
	if err := quic.StartService(config); err != nil {
		panic(err)
	}
}
