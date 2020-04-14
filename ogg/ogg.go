package ogg

import (
	"bytes"
	"fmt"
	"io"
)

type OggReader struct {
	reader io.Reader
	poz    int
	buffer []byte
}

func New(r io.Reader) *OggReader {
	return &OggReader{
		reader: r,
		poz:    -1,
		buffer: make([]byte, 0),
	}
}

func (r *OggReader) Page() (*Page, error) {
	for {
		for {
			if len(r.buffer) == 0 {
				break
			}
			poz := bytes.Index(r.buffer[r.poz+1:], []byte("OggS"))
			fmt.Println("poz", r.poz, poz, string(r.buffer))
			if poz == -1 {
				break
			}
			if r.poz == -1 {
				r.poz = poz
				continue
			}
			page := NewPage(r.buffer[r.poz : poz+1])
			r.poz = 0
			r.buffer = r.buffer[poz+1:]
			return page, nil
		}
		chunk := make([]byte, 500*1024)
		n, err := r.reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		r.buffer = append(r.buffer, chunk[:n]...)
	}
}
