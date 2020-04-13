package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/factorysh/streamcast/icecast"
)

type Streamer struct {
	clients map[int64]*flusherWriter
	id      int64
	lock    sync.Mutex
}

func (s *Streamer) nextID() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.id++
	return s.id
}

func New() *Streamer {
	return &Streamer{
		clients: make(map[int64]*flusherWriter),
	}
}

type flusherWriter struct {
	flusher http.Flusher
	writer  io.Writer
	end     context.CancelFunc
}

func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Proto)
	fmt.Println(r)
	/*if r.ProtoMajor != 2 {
		w.WriteHeader(http.StatusHTTPVersionNotSupported)
		return
	}*/
	switch r.Method {
	case "PUT":
		//	w.WriteHeader(100)
		defer r.Body.Close()
		for {
			var buff []byte
			n, err := r.Body.Read(buff)

			if err != nil {
				if err == io.EOF {
					fmt.Println("just an EOF.")
					time.Sleep(1000 * time.Millisecond)
				} else {
					fmt.Println("error", err)
					w.WriteHeader(500)
					return
				}
			}
			fmt.Println("buffer", n, buff)
		}
		w.WriteHeader(200)
	case "GET":
		f, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(500)
			return
		}
		ctx, cancelFunc := context.WithCancel(context.Background())
		s.clients[s.nextID()] = &flusherWriter{
			flusher: f,
			writer:  w,
			end:     cancelFunc,
		}
		select {
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	s := New()
	http.Handle("/", s)
	icecast.Listen("0.0.0.0:5000")
}
