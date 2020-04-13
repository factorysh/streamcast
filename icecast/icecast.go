package icecast

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func handleConnection(c net.Conn) {
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
		c.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n"))
		for {
			buff := make([]byte, 1024)
			n, err := b.Read(buff)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n == 0 {
				time.Sleep(100 * time.Millisecond)
			} else {
				fmt.Print(buff)
			}
		}
	}
}

func Listen(addr string) error {

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
		go handleConnection(c)
	}
}
