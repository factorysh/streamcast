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

func (s *Streamer) Read(buff []byte) {
	for _, client := range s.clients {
		client.writer.Write(buff)
		client.flusher.Flush()
	}
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
		w.Header().Add("Content-Type", "audio/ogg")
		w.Header().Add("Cache-Control", "no-cache")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Pragma", "no-cache")
		w.WriteHeader(200)
		ctx, cancelFunc := context.WithCancel(context.Background())
		id := s.nextID()
		s.clients[id] = &flusherWriter{
			flusher: f,
			writer:  w,
			end:     cancelFunc,
		}
		defer func() {
			delete(s.clients, id)
		}()
		select {
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	s := New()
	//h2s := &http2.Server{}
	mux := http.NewServeMux()
	mux.Handle("/stream", s)
	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	/*l, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		panic(err)
	}
	*/

	i := icecast.New()
	i.Reader = s.Read
	go i.Listen("0.0.0.0:5000")
	/*
		fmt.Println("Web")
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			h2s.ServeConn(conn, &http2.ServeConnOpts{
				Handler: mux,
			})
		}
	*/
	http.ListenAndServe(":5001", mux)
	/*
		h1s := &http.Server{
			Addr:    ":5001",
			Handler: h2c.NewHandler(mux, h2s),
		}
	*/
}
