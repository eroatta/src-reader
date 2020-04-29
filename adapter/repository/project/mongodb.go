package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoDBConnection = errors.New("MongoDB connection error")
	ErrMongoDBValidation = errors.New("MongoDB connection validation error")
)

type mongodb struct {
	db *mongo.Client
}

// NewMongoDB creates a repository.ProjectRepository backed up by a MongoDB database.
func NewMongoDB(client *mongo.Client) repository.ProjectRepository {
	return &mongodb{
		db: client,
	}
}

// NewMongoClient creates a new MongoDB client.
func NewMongoClient(url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.WithError(err).Error(fmt.Sprintf("error opening a connection to %s", url))
		return nil, ErrMongoDBConnection
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.WithError(err).Error(fmt.Sprintf("error validating a connection to %s", url))
		return nil, ErrMongoDBValidation
	}

	return client, nil
}

// TODO: improve
func (m *mongodb) Add(ctx context.Context, project entity.Project) error {
	projects := m.db.Database("reader").Collection("projects")
	res, err := projects.InsertOne(ctx, project)
	if err != nil {
		logrus.WithError(err).Error(fmt.Sprintf("error inserting record %v", project))
		return repository.ErrProjectUnexpected
	}
	logrus.WithField("name", project.Metadata.Fullname).Debug(fmt.Sprintf("inserted ID: %s", res.InsertedID))

	return nil
}

func (m *mongodb) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	return entity.Project{}, nil
}
