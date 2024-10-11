package coordinator

import (
	"WebRTC_POC/server/rtc/channel"
	"errors"
)

var (
	channelAlreadyExists = errors.New("channel already exists")
	channelNotExists     = errors.New("channel doesn't exist")
)

type Coordinator struct {
	channels map[string]*channel.Channel
}

func New() *Coordinator {
	return &Coordinator{
		channels: map[string]*channel.Channel{},
	}
}

func (s *Coordinator) CreateChannel(id string) (*channel.Channel, error) {
	ch, exist := s.channels[id]
	if !exist {
		ch = channel.New()
		s.channels[id] = ch
		return ch, nil
	} else {
		return ch, channelAlreadyExists
	}
}

func (s *Coordinator) GetChannel(id string) (*channel.Channel, error) {
	ch, exist := s.channels[id]
	if !exist {
		return nil, channelNotExists
	}
	return ch, nil
}

func (s *Coordinator) RemoveChannel(id string) {
	delete(s.channels, id)
}
