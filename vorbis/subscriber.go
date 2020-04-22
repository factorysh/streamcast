package vorbis

import (
	"context"
	"sync"

	"github.com/apex/log"
	"github.com/factorysh/streamcast/ogg"
)

type PubSub struct {
	streams    *Streams
	ventilator *ventilator
}

type ventilator struct {
	subscribers map[int64]*Subscriber
	cpt         int64
	lock        sync.RWMutex
}

func NewPubSub() *PubSub {
	v := &ventilator{
		subscribers: make(map[int64]*Subscriber),
	}
	s := NewStreams()
	s.Pipe = v
	return &PubSub{
		streams:    s,
		ventilator: v,
	}
}

func (v *ventilator) WritePage(page *ogg.Page) error {
	v.lock.RLock()
	defer v.lock.RUnlock()
	for _, subscriber := range v.subscribers {
		log.WithFields(log.Fields{
			"subscriber": subscriber,
			"serial":     page.Header().Serial,
			"granule":    page.Header().Granule,
		}).Info("Ventilator")
		subscriber.writer.Write(page.Raw)
		subscriber.writer.Flush()
	}
	// FIXME, never fail, really?
	return nil
}

func (p *PubSub) WritePage(page *ogg.Page) error {
	/*log.WithFields(log.Fields{
		"serial":  page.Header().Serial,
		"granule": page.Header().Granule,
	}).Info("Write a Page")*/
	return p.streams.WritePage(page)
}

type Subscriber struct {
	writer  WriterFlusher
	started bool
}

func (p *PubSub) Subscribe(ctx context.Context, w WriterFlusher) {
	p.ventilator.lock.Lock()
	id := p.ventilator.cpt
	sub := &Subscriber{
		writer:  w,
		started: false,
	}
	current := p.streams.CurrentStream()
	if current != nil {
		sub.started = current.WriteBegining(w)
	}
	p.ventilator.subscribers[id] = sub
	p.ventilator.cpt++
	p.ventilator.lock.Unlock()
	go func() {
		select {
		case <-ctx.Done():
			p.ventilator.lock.Lock()
			delete(p.ventilator.subscribers, id)
			p.ventilator.lock.Unlock()
		}
	}()
	log.Info("New subscriber")
}
