package tunnel

import (
	"bufio"
	"net"
)

type ServerTunnel struct {
}

// NewServer create a proxy server
func NewTunnelServer() *ServerTunnel {
	return &ServerTunnel{}
}

// newConn create a conn to serve client request
func (s *ServerTunnel) Process(rwc net.Conn, reader *bufio.Reader, requestTime int64) {
	run := &conn{
		server:      s,
		rwc:         rwc,
		brc:         reader,
		requestTime: requestTime,
	}
	run.serve()
}
