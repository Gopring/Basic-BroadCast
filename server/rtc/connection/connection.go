package connection

import (
	"WebRTC_POC/server/logging"
	"context"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v4"
)

type Connection struct {
	Conn *webrtc.PeerConnection
}

func NewInboundConnection(ctx context.Context, con webrtc.Configuration) (*Connection, error) {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}

	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}

	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}
	i.Add(intervalPliFactory)

	peerConnection, err := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i)).NewPeerConnection(con)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}
	return &Connection{Conn: peerConnection}, nil
}

func NewOutboundConnection(ctx context.Context, con webrtc.Configuration) (*Connection, error) {
	peerConnection, err := webrtc.NewPeerConnection(con)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}
	return &Connection{Conn: peerConnection}, nil
}

func (c *Connection) StartICE(ctx context.Context, sdp string) error {
	var err error
	broadOffer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}
	if err = c.Conn.SetRemoteDescription(broadOffer); err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}

	answer, err := c.Conn.CreateAnswer(nil)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}

	gatherComplete := webrtc.GatheringCompletePromise(c.Conn)

	err = c.Conn.SetLocalDescription(answer)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}
	<-gatherComplete
	return nil
}

func (c *Connection) ServerSDP() string {
	return c.Conn.LocalDescription().SDP
}
