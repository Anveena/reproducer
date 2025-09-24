package listener

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

type ioBufferSizeModel struct {
	readBufferSize  int
	writeBufferSize int
}

func ioBufferSize(c syscall.RawConn, bufferSize int, opt int) error {
	var rs error
	if err := c.Control(func(fd uintptr) {
		if rs = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, opt, bufferSize); rs != nil {
			rs = fmt.Errorf("SetsockoptInt failed: %s", rs.Error())
		}
	}); err != nil {
		if rs != nil {
			return fmt.Errorf("syscall.RawConn.Control failed: %s, %s", err.Error(), rs.Error())
		}
		return fmt.Errorf("syscall.RawConn.Control failed: %s", err.Error())
	}
	if rs != nil {
		return rs
	}
	if err := c.Control(func(fd uintptr) {
		actuallyValue, err := syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, opt)
		if err != nil {
			rs = fmt.Errorf("GetsockoptInt failed: %s", err.Error())
			return
		}
		if actuallyValue < bufferSize {
			rs = fmt.Errorf("set bufferSize failed, actually value: %d, needed value: %d", actuallyValue, bufferSize)
		}
	}); err != nil {
		if rs != nil {
			return fmt.Errorf("syscall.RawConn.Control failed: %s, %s", err.Error(), rs.Error())
		}
		return fmt.Errorf("syscall.RawConn.Control failed: %s", err.Error())
	}
	return rs
}

func newIOBufferSizeModel(readBufferSize, writeBufferSize int) *ioBufferSizeModel {
	return &ioBufferSizeModel{readBufferSize, writeBufferSize}
}

func (m *ioBufferSizeModel) listenerConfigControl(network, address string, c syscall.RawConn) error {
	if m.writeBufferSize > 0 {
		if err := ioBufferSize(c, m.writeBufferSize, syscall.SO_SNDBUF); err != nil {
			return fmt.Errorf("set SO_SNDBUF failed: %s", err.Error())
		}
	}
	if m.readBufferSize > 0 {
		if err := ioBufferSize(c, m.readBufferSize, syscall.SO_RCVBUF); err != nil {
			return fmt.Errorf("set SO_RCVBUF failed: %s", err.Error())
		}
	}
	return nil
}
func UDP(port uint16, readBufferSize int, writeBufferSize int) (*net.UDPConn, error) {
	bufferSize := newIOBufferSizeModel(readBufferSize, writeBufferSize)
	lc := net.ListenConfig{Control: bufferSize.listenerConfigControl}
	l, err := lc.ListenPacket(context.Background(), "udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("udp listen failed: %s", err.Error())
	}
	return l.(*net.UDPConn), nil
}
