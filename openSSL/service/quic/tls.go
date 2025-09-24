package quic

import (
	"crypto/tls"
	_ "embed"
)

var (
	//go:embed x509crt
	crtFile []byte
	//go:embed x509key
	keyFile []byte
)

var tlsConfig *tls.Config

func init() {
	certificate, err := tls.X509KeyPair(crtFile, keyFile)
	if err != nil {
		panic(err)
	}
	tlsConfig = &tls.Config{
		NextProtos:   []string{"anAlpnForTest"},
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
}
