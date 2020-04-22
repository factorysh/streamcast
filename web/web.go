package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/factorysh/streamcast/vorbis"
)

type Streamer struct {
	Pubsub *vorbis.PubSub
}

func NewStreamer() *Streamer {
	return &Streamer{
		Pubsub: vorbis.NewPubSub(),
	}
}

type ResponseWriterFlusher struct {
	writer  http.ResponseWriter
	flusher http.Flusher
}

func (r *ResponseWriterFlusher) Flush() {
	r.flusher.Flush()
}

func (r *ResponseWriterFlusher) Write(chunk []byte) {
	r.writer.Write(chunk)
}

func NewResponseWriterFlusher(w http.ResponseWriter) (*ResponseWriterFlusher, error) {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("Not flushable")
	}
	return &ResponseWriterFlusher{
		writer:  w,
		flusher: f,
	}, nil
}

func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "audio/ogg")
		w.Header().Add("Cache-Control", "no-cache")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Pragma", "no-cache")
		f, err := NewResponseWriterFlusher(w)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)

		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		s.Pubsub.Subscribe(ctx, f)
		select {
		case <-ctx.Done():
			log.Info("done")
		}
	}
}
