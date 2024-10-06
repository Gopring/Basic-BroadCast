package connection

import (
	"encoding/base64"
	"encoding/json"
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

func NewInboundConnection(con Config) (*Connection, error) {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		return nil, err
	}

	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		return nil, err
	}
	i.Add(intervalPliFactory)

	peerConnection, err := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i)).NewPeerConnection(con.peer)
	if err != nil {
		return nil, err
	}
	return &Connection{conn: peerConnection}, nil
}

func (c Connection) StartICE(sdp string) error {
	var err error
	broadOffer := webrtc.SessionDescription{}
	if err = decode(sdp, &broadOffer); err != nil {
		return err
	}

	if err = c.conn.SetRemoteDescription(broadOffer); err != nil {
		return err
	}

	answer, err := c.conn.CreateAnswer(nil)
	if err != nil {
		return err
	}

	gatherComplete := webrtc.GatheringCompletePromise(c.conn)

	err = c.conn.SetLocalDescription(answer)
	if err != nil {
		return err
	}
	<-gatherComplete
	return nil
}

func (c Connection) SetUpstream(stream *webrtc.TrackLocalStaticRTP) {
	c.conn.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}
			if _, err := stream.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				panic(err)
			}
		}
	})
}

func (c Connection) SetDownStream(stream *webrtc.TrackLocalStaticRTP) error {
	rtpSender, err := c.conn.AddTrack(stream)
	if err != nil {
		return err
	}
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
	return nil
}

func NewOutboundConnection(con Config) (*Connection, error) {
	peerConnection, err := webrtc.NewPeerConnection(con.peer)
	if err != nil {
		return nil, err
	}
	return &Connection{conn: peerConnection}, nil
}

func decode(in string, obj *webrtc.SessionDescription) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}
