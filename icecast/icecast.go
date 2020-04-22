package icecast

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/factorysh/streamcast/ogg"
)

type IcePUT struct {
	writer ogg.PageWriter
}

func New(writer ogg.PageWriter) *IcePUT {
	return &IcePUT{
		writer: writer,
	}
}

func (i *IcePUT) handleConnection(c net.Conn) {
	l := 0
	for {
		fmt.Println(c.RemoteAddr(), l, " ")
		b := bufio.NewReader(c)
		expect := false
		source := false
		for {
			netData, err := b.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Print(l, netData)
			if strings.HasPrefix(netData, "Expect") {
				expect = true
			}
			if strings.HasPrefix(netData, "SOURCE") && l == 0 {
				source = true
			}
			if netData == "\r\n" {
				break
			}
			l++
		}
		if expect {
			c.Write([]byte("HTTP/1.1 100 Continue\r\n"))
		}
		if source {
			c.Write([]byte("HTTP/1.0 200 OK\r\n"))
		}
		c.Write([]byte("Server: Streamcast\r\n"))
		c.Write([]byte("Connection: Close\r\n"))
		c.Write([]byte("Accept-Encoding: identity\r\n"))
		c.Write([]byte("Allow: GET, SOURCE\r\n"))
		c.Write([]byte("Cache-Control: no-cache\r\n"))
		c.Write([]byte("\r\n"))

		ctx := context.TODO()
		ogg.Stream(ctx, c, i.writer)
	}
}

func (i *IcePUT) Listen(addr string) error {

	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go i.handleConnection(c)
	}
}
