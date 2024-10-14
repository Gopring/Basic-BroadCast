package controller

import (
	"WebRTC_POC/server/backend"
	"WebRTC_POC/server/rtc/connection"
	"WebRTC_POC/types/request"
	"context"
	"errors"
)

var (
	InvalidRequest       = errors.New("invalid request data")
	ChannelAlreadyExists = errors.New("channel already exists")
	FailedNewConnection  = errors.New("failed to make connection")
	FailedNewStream      = errors.New("failed to make stream")
	FailedICE            = errors.New("failed to ICE")
	ChannelNotFound      = errors.New("channel not found")
)

type Controller struct {
	be *backend.Backend
}

func New(b *backend.Backend) *Controller {
	return &Controller{
		be: b,
	}
}

func (c *Controller) Broadcast(ctx context.Context) (string, error) {
	req := request.From(ctx)
	if req == nil {
		return "", InvalidRequest
	}

	ch, err := c.be.Coordinator.CreateChannel(req.ID)
	if err != nil {
		return "", ChannelAlreadyExists
	}

	conn, err := connection.NewInboundConnection(ctx, ch.Config)
	if err != nil {
		c.be.Coordinator.RemoveChannel(req.ID)
		return "", FailedNewConnection
	}

	ch.SetUpstream(ctx, conn, req.ID)

	err = conn.StartICE(ctx, req.SDP)
	if err != nil {
		c.be.Coordinator.RemoveChannel(req.ID)
		return "", FailedICE
	}
	return conn.ServerSDP(), nil
}

func (c *Controller) View(ctx context.Context) (string, error) {
	req := request.From(ctx)
	if req == nil {
		return "", InvalidRequest
	}

	ch, err := c.be.Coordinator.GetChannel(req.ID)
	if err != nil {
		return "", ChannelNotFound
	}

	conn, err := connection.NewOutboundConnection(ctx, ch.Config)
	if err != nil {
		return "", FailedNewConnection
	}

	err = ch.SetDownstream(ctx, conn)
	if err != nil {
		return "", FailedNewStream
	}

	err = conn.StartICE(ctx, req.SDP)
	if err != nil {
		return "", FailedICE
	}

	return conn.ServerSDP(), nil
}
