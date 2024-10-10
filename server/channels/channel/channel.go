package channel

import (
	"WebRTC_POC/server/channels/connection"
	"context"
	"github.com/pion/webrtc/v4"
)

type Channel struct {
	stream     *webrtc.TrackLocalStaticRTP
	connConfig connection.Config
}

func New() *Channel {
	return &Channel{}
}

func (c *Channel) SetBroadcaster(ctx context.Context, sdp string) error {
	conn, err := connection.NewInboundConnection(ctx, c.connConfig)
	if err != nil {
		return err
	}
	conn.SetUpstream(ctx, c.stream)
	return conn.StartICE(ctx, sdp)
}

func (c *Channel) SetViewer(ctx context.Context, sdp string) error {
	conn, err := connection.NewOutboundConnection(ctx, c.connConfig)
	if err != nil {
		return err
	}
	err = conn.SetDownStream(ctx, c.stream)
	return conn.StartICE(ctx, sdp)
}
