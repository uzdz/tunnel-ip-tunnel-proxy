package tunnel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"ip-tunnel-proxy/pkg/auth"
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/utils"
	"net"
	"net/textproto"
)

type conn struct {
	rwc         net.Conn
	brc         *bufio.Reader
	server      *ServerTunnel
	requestTime int64
}

func (c *conn) serve() {
	defer c.rwc.Close()

	rawHttpRequestHeader, credential, switchIp, err := c.getTunnelInfo()
	if err != nil {
		return
	}

	uid := auth.HttpAuthorization(credential, c.rwc.RemoteAddr().String())
	if config.NoAuth == 0 {
		if uid == "" {
			c.rwc.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"*\"\r\n\r\n"))
			return
		}
	}

	dstAddr, e := GetProxyDst(uid, switchIp)
	if e != nil {
		c.rwc.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n"))
		return
	}

	remoteConn, err := net.Dial("tcp", dstAddr)
	if err != nil {
		return
	}

	_, err = rawHttpRequestHeader.WriteTo(remoteConn)
	if err != nil {
		return
	}

	c.tunnel(remoteConn)
}

// getClientInfo parse client request header to get some information:
func (c *conn) getTunnelInfo() (rawReqHeader bytes.Buffer, credential string, switchIp bool, err error) {
	tp := textproto.NewReader(c.brc)

	// First line: GET /index.html HTTP/1.0
	var requestLine string
	if requestLine, err = tp.ReadLine(); err != nil {
		return
	}

	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return
	}

	credential = mimeHeader.Get(config.ProxyAuthorizationKey)

	if mimeHeader.Get(config.ProxySwitchIp) == "true" {
		switchIp = true
	}

	// rebuild http request header
	rawReqHeader.WriteString(requestLine + "\r\n")
	for k, vs := range mimeHeader {
		for _, v := range vs {
			if utils.In(k, config.IgnoreHeaderMap) {
				continue
			}
			rawReqHeader.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	rawReqHeader.WriteString(fmt.Sprintf("%s: %s\r\n", config.TunnelXForKey, config.TunnelXForValue))
	rawReqHeader.WriteString(fmt.Sprintf("%s: %s\r\n", config.XTunnelForwardedFor, c.rwc.RemoteAddr().String()))
	rawReqHeader.WriteString("\r\n")

	return
}

// tunnel http message between client and server
func (c *conn) tunnel(remoteConn net.Conn) {
	//inAddr := c.rwc.RemoteAddr().String()
	//inLocalAddr := c.rwc.LocalAddr().String()
	//
	//outAddr := remoteConn.RemoteAddr().String()
	//outLocalAddr := remoteConn.LocalAddr().String()

	//log.Printf("conn %s - %s - %s - %s connected", inAddr, inLocalAddr, outLocalAddr, outAddr)

	go func() {
		c.brc.WriteTo(remoteConn)
		//if err != nil {
		//	log.Println(err)
		//}
		remoteConn.Close()
	}()

	io.Copy(c.rwc, remoteConn)
	//if err != nil {
	//	log.Println(err)
	//}

	//log.Printf("conn %s - %s - %s -%s released", inAddr, inLocalAddr, outLocalAddr, outAddr)
}

type BadRequestError struct {
	what string
}

func (b *BadRequestError) Error() string {
	return b.what
}
