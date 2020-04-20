package ogg

import (
	"context"
	"sync"
)

type WriterFlusher interface {
	Flush()
	Write([]byte)
}

type Publisher struct {
	pageWriter PageWriter
	subscriber map[int64]WriterFlusher
	cpt        int64
	lock       sync.RWMutex
}

func NewPublisher(pw PageWriter) *Publisher {
	return &Publisher{
		pageWriter: pw,
		subscriber: make(map[int64]WriterFlusher),
	}
}

func (p *Publisher) Subscribe(ctx context.Context, w WriterFlusher) {
	p.lock.Lock()
	id := p.cpt
	p.subscriber[id] = w
	p.cpt++
	p.lock.Unlock()
	go func() {
		select {
		case <-ctx.Done():
			p.lock.Lock()
			delete(p.subscriber, id)
			p.lock.Unlock()
		}
	}()
}

func (p *Publisher) WriteAllSubscribers(chunk []byte) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, wf := range p.subscriber {
		wf.Write(chunk)
		wf.Flush()
	}
}
