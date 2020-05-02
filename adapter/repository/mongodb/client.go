package mongodb

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrMongoDBConnection represents an error while attempting to create a connection to the
	// given MongoDB server.
	ErrMongoDBConnection = errors.New("MongoDB connection error")
	// ErrorMongoDBValidation represents an error while attempting to check the existing connection
	// to the given MongoDB server.
	ErrMongoDBValidation = errors.New("MongoDB connection validation error")
)

// NewMongoClient creates a new MongoDB client.
func NewMongoClient(url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error opening a connection to %s", url))
		return nil, ErrMongoDBConnection
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error validating a connection to %s", url))
		return nil, ErrMongoDBValidation
	}

	return client, nil
}
