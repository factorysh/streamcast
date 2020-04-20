package vorbis

import (
	"errors"
	"fmt"

	"github.com/factorysh/streamcast/ogg"
)

type Streams struct {
	streams map[uint32]*Stream
}

func NewStreams() *Streams {
	return &Streams{
		streams: make(map[uint32]*Stream),
	}
}

func (s *Streams) WritePage(page *ogg.Page) error {
	header := page.Header()
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
		}
		s.streams[header.Serial] = stream
		return nil
	}
	if header.HeaderType == ogg.EOS {
		defer func() {
			delete(s.streams, header.Serial)
		}()
	}
	return stream.WritePage(page)
}

type Stream struct {
	headers []*ogg.Page
}

func NewStream() *Stream {
	return &Stream{
		headers: make([]*ogg.Page, 0),
	}
}

func (s *Stream) WritePage(page *ogg.Page) error {
	if len(s.headers) == 1 { // Vorbis starts with 2 pages of headers
		if page.Header().Granule != 0 {
			return errors.New("Vorbis starts with 2 pages of headers")
		}
		s.headers = append(s.headers, page)
		return nil
	}
	return nil
}

func (s *Stream) WriteBegining(w ogg.WriterFlusher) {
	// FIXME are all needed headers already here?
	for _, h := range s.headers {
		w.Write(h.Raw)
	}
	w.Flush()
}
