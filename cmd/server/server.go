package main

import (
	"context"

	"github.com/kubernetes-sigs/kernel-module-management/internal/proto"
)

type Server struct {
	proto.UnimplementedControlPlaneServer
}

func (s *Server) GetDesiredModules(ctx context.Context, node *proto.Node) (*proto.DesiredModules, error) {
	_ = node.GetName()

	dm := &proto.DesiredModules{
		Module: []*proto.Module{
			{
				K8SName:        "some-name",
				K8SNamespace:   "some-namespace",
				ContainerImage: "quay.io/path/test",
				Name:           "mymod",
				Parameters:     "",
			},
		},
	}

	return dm, nil
}
