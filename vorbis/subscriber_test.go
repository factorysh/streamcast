package vorbis

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/factorysh/streamcast/ogg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type mockWriter struct {
	buffer *bytes.Buffer
	cpt    int
}

func (m *mockWriter) Write(chunk []byte) {
	m.buffer.Write(chunk)
	m.cpt++
}

func (m *mockWriter) Flush() {}

func TestSubscriber(t *testing.T) {
	f, err := os.Open("../demo/slacker.ogg")
	assert.NoError(t, err)

	ctx := context.TODO()
	pubsub := NewPubSub()
	ctx2 := context.TODO()
	m2 := &mockWriter{
		buffer: &bytes.Buffer{},
	}
	pubsub.Subscribe(ctx2, m2)
	assert.Len(t, pubsub.ventilator.subscribers, 1)
	ogg.Stream(ctx, f, pubsub)
	time.Sleep(50 * time.Millisecond)
	assert.True(t, m2.cpt > 0)
	out := m2.buffer.Bytes()
	assert.True(t, len(out) > 0)
}
