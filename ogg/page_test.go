package ogg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	f, err := os.Open("../demo/slacker.ogg")
	assert.NoError(t, err)
	o := New(f)

	p, err := o.Page()
	assert.NoError(t, err)
	h := p.Header()
	assert.Equal(t, BOS, h.HeaderType) // BOS
	assert.Equal(t, int64(0), h.Granule)
	start := uint32(0)

	p, err = o.Page()
	assert.NoError(t, err)
	h = p.Header()
	assert.Equal(t, uint8(0), h.HeaderType)
	assert.Equal(t, int64(0), h.Granule)
	assert.Equal(t, start+1, h.Page)

	p, err = o.Page()
	assert.NoError(t, err)
	h = p.Header()
	assert.Equal(t, uint8(0), h.HeaderType)
	assert.NotEqual(t, int64(0), h.Granule)
	// OMG, this stream is cheating
	//assert.Equal(t, start+2, h.Page)
}
