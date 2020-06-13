package mongodb

import (
	"context"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
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

// Get finds an existing Project by ID.
func (pdb *ProjectDB) Get(ctx context.Context, ID uuid.UUID) (entity.Project, error) {
	return pdb.find(ctx, bson.M{"_id": ID.String})
}

// GetByReference finds an existing Project using the given reference as filter.
func (pdb *ProjectDB) GetByReference(ctx context.Context, projectRef string) (entity.Project, error) {
	return pdb.find(ctx, bson.M{"project_ref": projectRef})
}

func (pdb *ProjectDB) find(ctx context.Context, filter interface{}) (entity.Project, error) {
	res := pdb.collection.FindOne(ctx, filter)
	switch res.Err() {
	case nil:
		// do nothing
	case mongo.ErrNoDocuments:
		return entity.Project{}, repository.ErrProjectNoResults
	default:
		log.WithError(res.Err()).Errorf("error searching record with filter: %v", filter)
		return entity.Project{}, repository.ErrProjectUnexpected
	}

	var dto projectDTO
	if err := res.Decode(&dto); err != nil {
		log.WithError(err).Errorf("error decoding result for project with filter: %v", filter)
		return entity.Project{}, repository.ErrProjectUnexpected
	}

	return pdb.mapper.toEntity(dto), nil
}
