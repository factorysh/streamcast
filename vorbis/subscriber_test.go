package vorbis

import (
	"os"
	"testing"

	"github.com/factorysh/streamcast/ogg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSubscriber(t *testing.T) {
	f, err := os.Open("../demo/slacker.ogg")
	assert.NoError(t, err)

	ctx := context.TODO()
	pubsub := NewPubSub()
	ogg.Stream(ctx, f, pubsub)
}
