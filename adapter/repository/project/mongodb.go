package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoDBConnection = errors.New("MongoDB connection error")
	ErrMongoDBValidation = errors.New("MongoDB connection validation error")
)

type mongodb struct {
	client     *mongo.Client
	mapper     *projectMapper
	collection *mongo.Collection
}

// NewMongoDB creates a repository.ProjectRepository backed up by a MongoDB database.
func NewMongoDB(client *mongo.Client, dbname string) repository.ProjectRepository {
	return &mongodb{
		client:     client,
		mapper:     &projectMapper{},
		collection: client.Database(dbname).Collection("projects"),
	}
}

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

// Add transforms and stores a Project entity into a document on the underlying MongoDB collection.
func (m *mongodb) Add(ctx context.Context, project entity.Project) error {
	_, err := m.collection.InsertOne(ctx, m.mapper.toDTO(project))
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error inserting record %v", project))
		return repository.ErrProjectUnexpected
	}

	return nil
}

// GetByURL finds an existing Project using the given URL as filter.
func (m *mongodb) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	res := m.collection.FindOne(ctx, bson.M{"url": url})
	switch res.Err() {
	case nil:
		// do nothing
	case mongo.ErrNoDocuments:
		return entity.Project{}, repository.ErrProjectNoResults
	default:
		log.WithError(res.Err()).Error(fmt.Sprintf("error searching record with url: %s", url))
		return entity.Project{}, repository.ErrProjectUnexpected
	}

	var dto projectDTO
	if err := res.Decode(&dto); err != nil {
		log.WithError(err).Error(fmt.Sprintf("error decoding result for project with url: %s", url))
		return entity.Project{}, repository.ErrProjectUnexpected
	}

	return m.mapper.toEntity(dto), nil
}

func (m *mongodb) Close() {
	// TODO: close connections
}
