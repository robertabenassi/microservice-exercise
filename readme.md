# microservice-exercise
The microservice-exercise enables the update of a database of items.
Those items are given though a file which can be very big in size.

## Basic architecture
The microservice-exercise relies on a MongoDb database instance (so I am using a NoSQL solution for store data in a database). I will call it the mongodb service henceforth.
The service who is responsible to interact with the mongodb service is the port-domain-service, either creates new records in a database, or updates the existing ones.
The input needed for the database update is given through a file.
To present a possible way to upload a file, I am supposing that within the ecosystem in this exercise, there could be a web service with HTTP endpoint from which the upload is possible. 
The web service will be called port-api henceforth.

The microservice-exercise is a dockerized application, it has been developed/tested on:

## Docker Requirements
- docker (20.10.16)
- docker-compose (v1.29.3)

The docker-compose.yaml describes how services (the mongodb service, the port-domain-service, the port-api ) are connected and which are the dependencies among them.

## Additional Requirements
In order to provide a basic e2e test, we rely on curl which I consider the basic way to make HTTP requests (to the port-domain-service).
- curl (7.65.3)

## Protobuf
As far as microservices are concerned, I am supposing that the main communication paradigm should be grpc, so I am assuming that those tools are available
- protobuf 
- protoc-gen-go 
- protoc-gen-go-grpc 
(I had a bit of trouble on windows, things should go smoothier on Linux!)


## Useful commands on docker and utilities
In order to operate on microservices, I will put in a makefile the common used commands.
It is a very basic documentation!

## Protobuf generation
A Makefile is available on the folder proto, just to make easier the generation of protobuf go files.






