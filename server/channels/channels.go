package channels

import (
	"WebRTC_POC/server/channels/channel"
	"context"
	"fmt"
)

type Channels struct {
	channels map[string]*channel.Channel
}

func New() *Channels {
	return &Channels{
		channels: map[string]*channel.Channel{},
	}
}

func (s *Channels) Broadcast(ctx context.Context, key string, sdp string) error {
	ch, ok := s.channels[key]
	if !ok {
		ch = channel.New()
		s.channels[key] = ch
	}
	return ch.SetBroadcaster(ctx, sdp)
}

func (s *Channels) View(ctx context.Context, key string, sdp string) error {
	ch, ok := s.channels[key]
	if !ok {
		return fmt.Errorf("channel doesn't exist")
	}
	return ch.SetViewer(ctx, sdp)
}
