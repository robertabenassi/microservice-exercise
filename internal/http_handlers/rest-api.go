package http_handlers

import (
	"context"
	"io"
	"microservice-exercise/internal/grpc"
	"microservice-exercise/internal/stream"
	"microservice-exercise/internal/transport"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// RestClient is a client handling HTTP requests for ports
type RestClient struct {
	Client       transport.PortServiceClient
	parseService stream.StreamService
}

// NewRESTClient return new REST Client
func NewRESTClient(portServiceClient *grpc.Client, parseService stream.StreamService) *RestClient {
	return &RestClient{
		Client:       portServiceClient.Client,
		parseService: parseService,
	}
}

// HandleRequests, given a cancel function, manages the http requests (incoming requests)
func (c *RestClient) HandleRequests(cancel context.CancelFunc) {
	http.HandleFunc("/updatePorts", c.HandleUpdatePorts)
	log.Error(http.ListenAndServe(":5000", nil))
	cancel()
}

// HandleUpdatePorts handles HTTP requests to load ports from file (as a form load)
func (c *RestClient) HandleUpdatePorts(w http.ResponseWriter, r *http.Request) {
	log.Debug("Incoming requests updatePorts request")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Debug("Failed on loading file")
		log.Debug(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	f, err := os.OpenFile("/root/ports.json", os.O_WRONLY|os.O_CREATE, 0666) // what policy for file permission?
	if err != nil {
		log.Debug("Failed on opening file")
		log.Debug(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	io.Copy(f, file)
	go c.parseService.Load(f.Name())
	w.WriteHeader(http.StatusOK)
}
