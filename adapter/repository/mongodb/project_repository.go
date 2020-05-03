package mongodb

import (
	"context"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const projectsCollection string = "projects"

// ProjectDB represents a MongoDB database, focused on the collection handling the project documents.
type ProjectDB struct {
	client     *mongo.Client
	mapper     *projectMapper
	collection *mongo.Collection
}

// NewMongoDBProjecRepository creates a repository.ProjectRepository backed up by a MongoDB database.
func NewMongoDBProjecRepository(client *mongo.Client, dbname string) *ProjectDB {
	return &ProjectDB{
		client:     client,
		mapper:     &projectMapper{},
		collection: client.Database(dbname).Collection(projectsCollection),
	}
}

// Add transforms and stores a Project entity into a document on the underlying MongoDB collection.
func (pdb *ProjectDB) Add(ctx context.Context, project entity.Project) error {
	_, err := pdb.collection.InsertOne(ctx, pdb.mapper.toDTO(project))
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error inserting record %v", project))
		return repository.ErrProjectUnexpected
	}

	return nil
}

// GetByURL finds an existing Project using the given URL as filter.
func (pdb *ProjectDB) GetByURL(ctx context.Context, url string) (entity.Project, error) {
	res := pdb.collection.FindOne(ctx, bson.M{"url": url})
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

	return pdb.mapper.toEntity(dto), nil
}
