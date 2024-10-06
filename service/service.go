package service

import (
	"WebRTC_POC/channel"
	"fmt"
)

type Service struct {
	channels map[string]*channel.Channel
}

func (s *Service) Broadcast(key string, sdp string) error {
	ch, ok := s.channels[key]
	if !ok {
		s.channels[key] = channel.NewChannel()
	}
	return ch.SetBroadcaster(sdp)
}

func (s *Service) View(key string, sdp string) error {
	ch, ok := s.channels[key]
	if !ok {
		return fmt.Errorf("channel doesn't exist")
	}
	return ch.SetBroadcaster(sdp)
}
