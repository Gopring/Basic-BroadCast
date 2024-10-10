package connection

import (
	"WebRTC_POC/server/logging"
	"context"
	"errors"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v4"
	"io"
)

type Connection struct {
	conn *webrtc.PeerConnection
}

type Config struct {
	peer webrtc.Configuration
}

func NewInboundConnection(ctx context.Context, con Config) (*Connection, error) {
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

	peerConnection, err := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i)).NewPeerConnection(con.peer)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}
	return &Connection{conn: peerConnection}, nil
}

func (c Connection) StartICE(ctx context.Context, sdp string) error {
	var err error
	broadOffer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}

	if err = c.conn.SetRemoteDescription(broadOffer); err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}

	answer, err := c.conn.CreateAnswer(nil)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}

	gatherComplete := webrtc.GatheringCompletePromise(c.conn)

	err = c.conn.SetLocalDescription(answer)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}
	<-gatherComplete
	return nil
}

func (c Connection) SetUpstream(ctx context.Context, stream *webrtc.TrackLocalStaticRTP) {
	logging.DefaultLogger().Debugf(`upstreamID: %s`, stream.StreamID())
	c.conn.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				logging.From(ctx).Named("connection").Panic(readErr)
				panic(readErr)
			}
			if _, err := stream.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				logging.From(ctx).Named("connection").Panic(readErr)
				panic(err)
			}
		}
	})
}

func (c Connection) SetDownStream(ctx context.Context, stream *webrtc.TrackLocalStaticRTP) error {
	logging.DefaultLogger().Debugf(`downstreamID: %s`, stream.StreamID())
	rtpSender, err := c.conn.AddTrack(stream)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return err
	}
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				logging.From(ctx).Named("connection").Error(rtcpErr)
				return
			}
		}
	}()
	return nil
}

func NewOutboundConnection(ctx context.Context, con Config) (*Connection, error) {
	peerConnection, err := webrtc.NewPeerConnection(con.peer)
	if err != nil {
		logging.From(ctx).Named("connection").Error(err)
		return nil, err
	}
	return &Connection{conn: peerConnection}, nil
}
