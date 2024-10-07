package coordinator

import (
	"WebRTC_POC/pkg/channel"
	"fmt"
)

type Coordinator struct {
	channels map[string]*channel.Channel
}

func New() *Coordinator {
	return &Coordinator{
		channels: map[string]*channel.Channel{},
	}
}

func (s *Coordinator) Broadcast(key string, sdp string) error {
	ch, ok := s.channels[key]
	if !ok {
		s.channels[key] = channel.New()
	}
	return ch.SetBroadcaster(sdp)
}

func (s *Coordinator) View(key string, sdp string) error {
	ch, ok := s.channels[key]
	if !ok {
		return fmt.Errorf("channel doesn't exist")
	}
	return ch.SetBroadcaster(sdp)
}
