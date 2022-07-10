package grpc

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"microservice-exercise/internal/data_model"
	"microservice-exercise/internal/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client is a client handling gRPC requests
type Client struct {
	Client     transport.PortServiceClient
	connection *grpc.ClientConn
}

// NewClient returns a new Client
func NewClient(serverAddr string, cancel context.CancelFunc) *Client {
	var opts []grpc.DialOption
	// insecure service!!! what should we do about that?
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Info("Dialing ", serverAddr, "...")
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Error("failed to dial: ", err)
		cancel()
	}
	client := transport.NewPortServiceClient(conn)
	return &Client{Client: client, connection: conn}
}

// UpdatePort handles grpc requests for creating or upating port item
func (c *Client) UpdatePort(ID string, port data_model.Port) error {
	log.Println("UpdatePort ", ID, " port")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.Client.UpdatePort(ctx, &data_model.UpdatePortRequest{
		Key:  ID,
		Port: &port,
	})
	return err
}

// CloseConnection closes the client connections
func (c *Client) CloseConnection() error {
	return c.connection.Close()
}
