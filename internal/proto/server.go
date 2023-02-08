package proto

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/metadata"
)

type ServerImpl struct {
	logger logr.Logger
}

func NewServerImpl(logger logr.Logger) *ServerImpl {
	return &ServerImpl{logger: logger}
}

func (s *ServerImpl) GetDaemonConfig(ctx context.Context, request *DaemonConfigRequest) (*DaemonConfigResponse, error) {
	md, found := metadata.FromIncomingContext(ctx)

	s.logger.Info(
		"New request",
		"metadata", md,
		"found", found,
		"hostname", request.HostName,
		"kernelVersion", request.KernelVersion,
	)

	select {
	case <-time.After(5 * time.Second):
		s.logger.Info("Work completed")
	case <-ctx.Done():
		s.logger.Info("Cancelled by request context")
	}

	return nil, errors.New("not implemented")
}

func (s *ServerImpl) mustEmbedUnimplementedKMMServer() {
	//TODO implement me
	panic("implement me")
}
