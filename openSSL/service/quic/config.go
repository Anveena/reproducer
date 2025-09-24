package quic

type Config struct {
	Port                   uint16
	ReadChanSizeForClient  uint
	WriteChanSizeForClient uint
	ReadBufferSize         int
	WriteBufferSize        int
}
