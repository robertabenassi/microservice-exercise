package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"microservice-exercise/internal/database"
	"microservice-exercise/internal/grpc"
)

type PortDatabaseServiceConfiguration struct {
	Port    int
	Address string
}

func main() {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-sigCh
		log.Error("received a system call: ", oscall)
		cancel()
	}()

	configuration, err := getDatabaseServerParameters()
	if err != nil {
		log.Error(err)
		cancel()
	}

	mongodbClient, err := database.NewMongoClient(configuration.Address)
	if err != nil {
		log.Error(err)
		cancel()
	}
	defer mongodbClient.Client.Disconnect(ctx)
	grpcServer := grpc.NewPortsProtocolServer(configuration.Port, mongodbClient, cancel)
	defer grpcServer.GracefulStop()
	go func() {
		if err = grpcServer.Serve(); err != nil {
			log.Error(err)
			cancel()
		}
	}()
	<-ctx.Done()
	log.Info("Stopping Port Domain Service")
}

func getDatabaseServerParameters() (PortDatabaseServiceConfiguration, error) {
	var configuration PortDatabaseServiceConfiguration

	port := flag.Int("port", 5000, "Server port")
	address := flag.String("mongodb-address", "mongodb_exercise:27017", "Address of mongoDB database instance")
	logLevel := flag.String("log-level", "info", "Log level.")
	flag.Parse()
	logLevelParsed, err := log.ParseLevel(*logLevel)
	if err != nil {
		return configuration, err
	}
	log.SetLevel(logLevelParsed)

	configuration.Port = *port
	configuration.Address = *address

	log.Debug("Running with CLI parameters. port: ", configuration.Port,
		" mongodb-address: ", configuration.Address,
		" log-level: ", *logLevel)
	return configuration, nil
}
