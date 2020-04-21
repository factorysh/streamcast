package vorbis

import (
	"bytes"
	"os"
	"testing"

	"github.com/factorysh/streamcast/ogg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type mockWriter struct {
	buffer *bytes.Buffer
}

func (m *mockWriter) Write(chunk []byte) {
	m.buffer.Write(chunk)
}

func (m *mockWriter) Flush() {}

func TestSubscriber(t *testing.T) {
	f, err := os.Open("../demo/slacker.ogg")
	assert.NoError(t, err)

	ctx := context.TODO()
	pubsub := NewPubSub()
	ogg.Stream(ctx, f, pubsub)
	ctx2 := context.TODO()
	m2 := &mockWriter{
		buffer: &bytes.Buffer{},
	}
	pubsub.Subscribe(ctx2, m2)
}
