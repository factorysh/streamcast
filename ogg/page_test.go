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
	assert.Equal(t, uint8(2), h.HeaderType) // BOS
	assert.Equal(t, int64(0), h.Granule)
	p, err = o.Page()
	assert.NoError(t, err)
	h = p.Header()
	assert.Equal(t, uint8(0), h.HeaderType)
	assert.Equal(t, int64(0), h.Granule)
	p, err = o.Page()
	assert.NoError(t, err)
	h = p.Header()
	assert.Equal(t, uint8(0), h.HeaderType)
	assert.NotEqual(t, int64(0), h.Granule)
}
