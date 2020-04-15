package ogg

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOgg(t *testing.T) {
	r := bytes.NewReader([]byte(`OggS_beuha_OggS_aussi_OggS`))
	fmt.Println(r.Len())
	o := New(r)
	p, err := o.Page()
	assert.NoError(t, err)
	assert.Equal(t, "OggS_beuha_", string(p.Raw))
	p, err = o.Page()
	assert.NoError(t, err)
	assert.Equal(t, "OggS_aussi_", string(p.Raw))
}
