package vorbis

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/factorysh/streamcast/ogg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestStream(t *testing.T) {
	f, err := os.Open("../demo/slacker.ogg")
	assert.NoError(t, err)

	ctx := context.TODO()
	streams := NewStreams()
	err = ogg.Stream(ctx, f, streams)
	fmt.Println(streams.streams)
	assert.Len(t, streams.streams, 1)
	assert.NotNil(t, err)
	assert.Equal(t, io.EOF, err)
}
