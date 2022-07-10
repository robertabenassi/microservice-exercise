package faker

import (
	"context"
	"microservice-exercise/internal/data_model"
	"microservice-exercise/internal/transport"

	"google.golang.org/grpc"
)

type fakePortServiceClient struct {
}

// NewFakePortServiceClient returns fakePortServiceClient interface implementation
func NewFakePortServiceClient() transport.PortServiceClient {
	return &fakePortServiceClient{}
}

// UpdatePort is an implementation of the PortServiceClient interface
func (c *fakePortServiceClient) UpdatePort(ctx context.Context, in *data_model.UpdatePortRequest, opts ...grpc.CallOption) (*data_model.UpdatePortResponse, error) {
	return &data_model.UpdatePortResponse{}, nil
}
