package database

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus" // it is the most common used library for logging,

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"microservice-exercise/internal/data_model"
)

type portDB struct {
	ID   string          `bson:"_id"`
	Port data_model.Port `bson:"port"`
}

// MongoClient is a struct for representing a MongoDb Client.
type MongoClient struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// NewMongoClient returns a mongoDB client, opening the connection to the provided address.
// It returns an error if the connection is not working (i.e. the application should not start without the service)
func NewMongoClient(address string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%v", address))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return &MongoClient{}, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error("Got an error while tring to ping the MongoDB datavase, error was ", err)
		return &MongoClient{}, err
	}
	log.Info("Just connected to MongoDB at ", address)
	collection := client.Database("port").Collection("port")

	return &MongoClient{
		Client:     client,
		Collection: collection,
	}, nil
}

// UpdatePort updates a port in the port collections
func (m *MongoClient) UpdatePort(PortID string, port data_model.Port) error {
	filter := bson.M{"_id": PortID}
	portEntry := portDB{
		ID:   PortID,
		Port: port,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.Collection.FindOne(ctx, filter).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			insertResult, err := m.Collection.InsertOne(context.TODO(), portEntry)
			log.Debug("Insert operation for ", PortID, " got the result: ", insertResult)
			if err != nil {
				log.Error("Insert operation for ", PortID, " got an error: ", err)
			}
			return err
		}
	}
	replaceResult, err := m.Collection.ReplaceOne(context.TODO(), filter, portEntry)
	log.Debug("Replace operation for ", PortID, " got the result: ", replaceResult)
	if err != nil {
		log.Error("Insert operation for ", PortID, " got an error: ", err)
	}
	return err
}
