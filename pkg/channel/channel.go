package channel

import (
	"WebRTC_POC/pkg/connection"
	"github.com/pion/webrtc/v4"
)

type Channel struct {
	stream     *webrtc.TrackLocalStaticRTP
	connConfig connection.Config
}

func New() *Channel {
	return &Channel{}
}

func (c Channel) SetBroadcaster(sdp string) error {
	conn, err := connection.NewInboundConnection(c.connConfig)
	if err != nil {
		return err
	}
	conn.SetUpstream(c.stream)
	return conn.StartICE(sdp)
}

func (c Channel) SetViewer(sdp string) error {
	conn, err := connection.NewOutboundConnection(c.connConfig)
	if err != nil {
		return err
	}
	err = conn.SetDownStream(c.stream)
	return conn.StartICE(sdp)
}
