package server

import (
	"bufio"
	"ip-tunnel-proxy/pkg/tunnel"
	"log"
	"net"
	"strings"
	"time"
)

var P *S

type Server interface {
	Process(conn net.Conn, reader *bufio.Reader, requestTime int64)
}

type S struct {
	listener net.Listener
}

func (s *S) OpenServer(port string) {

	tunnelServer := tunnel.NewTunnelServer()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}
	s.listener = listener

	log.Println("代理服务器启动成功...")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("正在关闭代理服务器...")
				break
			}
			log.Println(err)
			continue
		}

		requestTime := time.Now().Unix()
		bufConn := bufio.NewReader(conn)
		go tunnelServer.Process(conn, bufConn, requestTime)
	}
}

func (s *S) Shutdown() bool {
	err := s.listener.Close()

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
