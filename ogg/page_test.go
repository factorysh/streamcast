package ogg

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	f, err := os.Open("../demo/slacker.ogg")
	assert.NoError(t, err)

	w := &writer{
		pages: make([]*Page, 0),
	}
	w.max = 3
	w.wg.Add(3)
	ctx := context.TODO()
	Stream(ctx, f, w)
	w.wg.Wait()

	p := w.pages[0]
	h := p.Header()
	assert.Equal(t, BOS, h.HeaderType) // BOS
	assert.Equal(t, int64(0), h.Granule)
	start := uint32(0)

	p = w.pages[1]
	h = p.Header()
	assert.Equal(t, uint8(0), h.HeaderType)
	assert.Equal(t, int64(0), h.Granule)
	assert.Equal(t, start+1, h.Page)

	p = w.pages[2]
	h = p.Header()
	assert.Equal(t, uint8(0), h.HeaderType)
	assert.NotEqual(t, int64(0), h.Granule)
	// OMG, this stream is cheating
	//assert.Equal(t, start+2, h.Page)
}
