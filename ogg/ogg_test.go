package ogg

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type writer struct {
	pages []*Page
	wg    sync.WaitGroup
	max   int
}

func (w *writer) WritePage(p *Page) error {
	w.pages = append(w.pages, p)
	w.max--
	if w.max >= 0 {
		w.wg.Done()

	}
	return nil
}

func TestOgg(t *testing.T) {
	r := bytes.NewReader([]byte(`OggS_beuha_OggS_aussi_OggS`))
	fmt.Println(r.Len())
	ctx := context.TODO()
	w := &writer{
		pages: make([]*Page, 0),
	}
	w.max = 2
	w.wg.Add(2)
	go Stream(ctx, r, w)
	w.wg.Wait()

	assert.Equal(t, "OggS_beuha_", string(w.pages[0].Raw))
	assert.Equal(t, "OggS_aussi_", string(w.pages[1].Raw))
}
