package vorbis

import (
	"errors"
	"fmt"
	"sync"

	"github.com/factorysh/streamcast/ogg"
)

type WriterFlusher interface {
	Flush()
	Write([]byte)
}

type Streams struct {
	streams map[uint32]*Stream
	current uint32
	lock    sync.RWMutex
	Pipe    ogg.PageWriter
}

func NewStreams() *Streams {
	return &Streams{
		streams: make(map[uint32]*Stream),
	}
}

func (s *Streams) CurrentStream() *Stream {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.streams[s.current]
}

func (s *Streams) WritePage(page *ogg.Page) error {
	header := page.Header()
	s.lock.Lock()
	defer s.lock.Unlock()
	stream, ok := s.streams[header.Serial]
	if !ok { // new Stream
		if header.HeaderType != ogg.BOS {
			return fmt.Errorf("Bad type, it must be begining, not %v", header.HeaderType)
		}
		if header.Granule != 0 {
			return fmt.Errorf("Screwed vorbis page, it must begin with granule 0, not %v", header.Granule)
		}
		stream = &Stream{
			headers: []*ogg.Page{page},
			Pipe:    s.Pipe,
		}
		s.streams[header.Serial] = stream
		s.current = header.Serial
		return nil
	}
	if header.HeaderType == ogg.EOS {
		s.lock.Unlock()
		defer func() {
			s.lock.Lock()
			delete(s.streams, header.Serial)
			s.lock.Unlock()
		}()
	}
	return stream.WritePage(page)
}

type Stream struct {
	headers []*ogg.Page
	Pipe    ogg.PageWriter
}

func NewStream() *Stream {
	return &Stream{
		headers: make([]*ogg.Page, 0),
	}
}

func (s *Stream) WritePage(page *ogg.Page) error {
	if s.Pipe != nil {
		s.Pipe.WritePage(page)
	}
	if len(s.headers) == 1 { // Vorbis starts with 2 pages of headers
		if page.Header().Granule != 0 {
			return errors.New("Vorbis starts with 2 pages of headers")
		}
		s.headers = append(s.headers, page)
		return nil
	}
	return nil
}

func (s *Stream) WriteBegining(w WriterFlusher) bool {
	// FIXME are all needed headers already here?
	if len(s.headers) < 2 {
		return false
	}
	for _, h := range s.headers {
		w.Write(h.Raw)
	}
	w.Flush()
	return true
}
