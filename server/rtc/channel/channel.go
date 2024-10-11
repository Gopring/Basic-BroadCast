package channel

import (
	"WebRTC_POC/server/logging"
	"WebRTC_POC/server/rtc/connection"
	"context"
	"errors"
	"github.com/pion/webrtc/v4"
	"io"
)

type Channel struct {
	stream *webrtc.TrackLocalStaticRTP
	Config webrtc.Configuration
}

func New() *Channel {
	return &Channel{}
}

func (c *Channel) SetUpstream(ctx context.Context, conn *connection.Connection, id string) {
	conn.Conn.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		var newTrackErr error
		c.stream, newTrackErr = webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "video", id)
		if newTrackErr != nil {
			logging.From(ctx).Named("connection").Panic(newTrackErr)
			panic(newTrackErr)
		}

		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				logging.From(ctx).Named("connection").Error(readErr)
				panic(readErr)
			}
			if _, err := c.stream.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				logging.From(ctx).Named("connection").Error(readErr)
				panic(err)
			}
		}
	})
}

func (c *Channel) SetDownstream(ctx context.Context, conn *connection.Connection) error {
	if c.stream == nil {
		logging.From(ctx).Named("connection").Error("stream not exists")
		return errors.New("stream not exists")
	}

	rtpSender, err := conn.Conn.AddTrack(c.stream)
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
