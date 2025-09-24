package quic

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/Anveena/opensslreproducer/listener"
	"github.com/quic-go/quic-go"
)

func StartService(config *Config) error {
	runtime.LockOSThread()
	quicConfig := new(quic.Config)
	// Allow0RTT,EnableDatagrams Listen always set.
	quicConfig.Allow0RTT = true
	quicConfig.EnableDatagrams = true
	quicConfig.MaxIdleTimeout = time.Hour * 24 * 365
	quicConfig.KeepAlivePeriod = time.Second * 30
	udpConn, err := listener.UDP(config.Port, config.ReadBufferSize, config.WriteBufferSize)
	if err != nil {
		return fmt.Errorf("udp listen failed: %s", err.Error())
	}
	quicListener, err := quic.Listen(udpConn,
		tlsConfig,
		quicConfig)
	if err != nil {
		return fmt.Errorf("quic listen failed: %s", err.Error())
	}
	ctx := context.Background()
	for {
		conn, err := quicListener.Accept(ctx)
		if err != nil {
			return fmt.Errorf("quic accept failed: %s", err.Error())
		}
		go onNewQUICConn(conn)
	}
}
func onNewQUICConn(conn *quic.Conn) {
	defer func() {
		_ = conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
		println("closed client")
	}()
	println("accepting stream")
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return
	}
	defer func() {
		_ = stream.Close()
		println("closed stream")
	}()
	buffer := make([]byte, 8)
	if _, err := io.ReadFull(stream, buffer); err != nil {
		println(err.Error())
		return
	}
	generation := binary.LittleEndian.Uint32(buffer)
	total := binary.LittleEndian.Uint32(buffer[4:])
	rsp := make([]byte, 256)
	for i := range 256 {
		rsp[i] = byte(i)
	}
	_, err = stream.Write(rsp)
	if err != nil {
		println(err.Error())
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	getManager().join(generation, total, &wg)
	wg.Wait()
}
