package icecast

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/mccoyst/ogg"
)

type IcePUT struct {
	Stream chan []byte
	Reader func([]byte)
}

func New() *IcePUT {
	return &IcePUT{
		Stream: make(chan []byte),
		Reader: func(buff []byte) {
			fmt.Print(len(buff))
		},
	}
}

func (i *IcePUT) handleConnection(c net.Conn) {
	l := 0
	for {
		fmt.Print(c.RemoteAddr(), l, " ")
		b := bufio.NewReader(c)
		for {
			netData, err := b.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Print(netData)
			if netData == "\r\n" {
				break
			}
			l++
		}
		c.Write([]byte("HTTP/1.1 100 Continue\r\n"))
		c.Write([]byte("Server: Streamcast\r\n"))
		c.Write([]byte("Connection: Close\r\n"))
		c.Write([]byte("Accept-Encoding: identity\r\n"))
		c.Write([]byte("Allow: GET, SOURCE\r\n"))
		c.Write([]byte("Cache-Control: no-cache\r\n"))
		c.Write([]byte("\r\n"))

		var buffer *bytes.Buffer
		for {
			chunk := make([]byte, 1024*1024)
			n, err := b.Read(chunk)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if buffer == nil {
				buffer = bytes.NewBuffer(chunk)
			} else {
				buffer.Write(chunk)
			}
			for {
				decoder := ogg.NewDecoder(buffer)
				page, err := decoder.Decode()
				if err != nil {
					if err == io.ErrUnexpectedEOF || err == io.EOF {
						break
					}
					fmt.Println(err)
					return
				}
				fmt.Println(page.Type)
				raw := page.Packet
				i.Reader(raw)
			}
			//fmt.Println(n)
		}
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
