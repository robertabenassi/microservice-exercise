package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"microservice-exercise/internal/grpc"
	"microservice-exercise/internal/http_handlers"
	"microservice-exercise/internal/stream"
)

type PortAPIConfiguration struct {
	Port    int
	Address string
}

func main() {
	// instance a channel trough which listen for syscall
	sigCh := make(chan os.Signal, 1)
	// the client should gracefully stops, expecially on SIGTERM
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// get a context and a cancel function to handle the stop request gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		oscall := <-sigCh
		log.Error("received a system call: ", oscall)
		cancel()
	}()

	configuration, err := getPortAPIConfiguration()
	if err != nil {
		log.Error(err)
		cancel()
	}
	// formatting a url of the portAPI
	serverAddr := fmt.Sprintf("%v:%d", configuration.Address, configuration.Port)

	stream := stream.NewStream()
	Client := grpc.NewClient(serverAddr, cancel)
	defer Client.CloseConnection() // ensure that everything is closed before exiting/stopping

	// a go routine for handling the JSON file of ports as a stream
	go watchJSONPortStream(cancel, stream, Client)

	// the http server to handle requests as a go routine
	go http_handlers.NewRESTClient(Client, stream).HandleRequests(cancel)

	<-ctx.Done()
	log.Info("Stopping reading the JSON file stream")
}

func watchJSONPortStream(cancel context.CancelFunc, stream *stream.Stream, Client *grpc.Client) error {
	for data := range stream.Watch() {
		if data.Error != nil {
			log.Error(data.Error)
			cancel()
		}
		log.Info("creating ", data.ID, " : ", data.Port.Name)
		if err := Client.UpdatePort(data.ID, data.Port); err != nil {
			log.Error(err)
			cancel()
		}
	}
	return nil
}

func getPortAPIConfiguration() (PortAPIConfiguration, error) {
	var configuration PortAPIConfiguration

	port := flag.Int("port-domain-server", 5000, "Domain server port.")
	address := flag.String("server-address", "portdomainservice_exercise", "Address of the port domain server.")
	logLevel := flag.String("log-level", "info", "Log level.")
	flag.Parse()
	logLevelParsed, err := log.ParseLevel(*logLevel)
	if err != nil {
		return configuration, err
	}
	log.SetLevel(logLevelParsed)

	configuration.Port = *port
	configuration.Address = *address

	log.Info("Running with CLI parameters. port-domain-server: ", configuration.Port,
		" server-address: ", configuration.Address,
		" log-level: ", *logLevel)
	return configuration, nil
}
