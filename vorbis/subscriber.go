package vorbis

import (
	"context"
	"sync"

	"github.com/factorysh/streamcast/ogg"
)

type PubSub struct {
	streams    *Streams
	ventilator *ventilator
}

type ventilator struct {
	subscribers map[int64]ogg.WriterFlusher
	cpt         int64
	lock        sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		streams: NewStreams(),
		ventilator: &ventilator{
			subscribers: make(map[int64]ogg.WriterFlusher),
		},
	}
}

func (v *ventilator) WritePage(page *ogg.Page) error {
	v.lock.RLock()
	defer v.lock.Unlock()
	for _, subscriber := range v.subscribers {
		subscriber.Write(page.Raw)
		subscriber.Flush()
	}
	// FIXME, never fail, really?
	return nil
}

func (p *PubSub) WritePage(page *ogg.Page) error {
	return p.streams.WritePage(page)
}

func (p *PubSub) Subscribe(ctx context.Context, w ogg.WriterFlusher) {
	p.ventilator.lock.Lock()
	id := p.ventilator.cpt
	p.ventilator.subscribers[id] = w
	p.ventilator.cpt++
	p.streams.CurrentStream().WriteBegining(w)
	p.ventilator.lock.Unlock()
	go func() {
		select {
		case <-ctx.Done():
			p.ventilator.lock.Lock()
			delete(p.ventilator.subscribers, id)
			p.ventilator.lock.Unlock()
		}
	}()
}
