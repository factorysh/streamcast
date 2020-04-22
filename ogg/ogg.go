package ogg

import (
	"bytes"
	"context"
	"io"
)

type PageWriter interface {
	WritePage(*Page) error
}

type OggReader struct {
	reader io.Reader
	writer PageWriter
	poz    int
	buffer []byte
	errors chan error
}

// Stream read a io.Reader and writes Page
func Stream(ctx context.Context, r io.Reader, w PageWriter) error {
	o := &OggReader{
		reader: r,
		writer: w,
		poz:    -1,
		buffer: make([]byte, 0),
		errors: make(chan error),
	}
	go func() {
		for {
			o.streamOnePage()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case error := <-o.errors:
			return error
		}
	}
}

func (r *OggReader) streamOnePage() {
	for {
		for {
			if len(r.buffer) == 0 {
				break
			}
			poz := bytes.Index(r.buffer[r.poz+1:], []byte("OggS"))
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
			err := r.writer.WritePage(page)
			if err != nil {
				r.errors <- err
			}
			return
		}
		chunk := make([]byte, 500*1024)
		n, err := r.reader.Read(chunk)
		if err != nil {
			r.errors <- err
			return
		}
		r.buffer = append(r.buffer, chunk[:n]...)
	}
}
