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

const analysisCollection string = "analysis"

// AnalysisDB represents a MongoDB database, focused on the collection handling the analysis documents.
type AnalysisDB struct {
	client     *mongo.Client
	mapper     *analysisMapper
	collection *mongo.Collection
}

// NewMongoDBAnalysisRepository creates a repository.AnalysisRepository backed up by a MongoDB database.
func NewMongoDBAnalysisRepository(client *mongo.Client, dbname string) *AnalysisDB {
	return &AnalysisDB{
		client:     client,
		mapper:     &analysisMapper{},
		collection: client.Database(dbname).Collection(analysisCollection),
	}
}

// Add transforms and stores an AnalysisResults entity into a document on the underlying MongoDB collection.
func (adb *AnalysisDB) Add(ctx context.Context, analysis entity.AnalysisResults) error {
	_, err := adb.collection.InsertOne(ctx, adb.mapper.toDTO(analysis))
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error inserting record %v", analysis))
		return repository.ErrAnalysisUnexpected
	}

	return nil
}

// GetByProjectID retrieves an existing analysis for the given Project, from the underlying MongoDB collection.
func (adb *AnalysisDB) GetByProjectID(ctx context.Context, projectID uuid.UUID) (entity.AnalysisResults, error) {
	results := adb.collection.FindOne(ctx, bson.M{"project_id": projectID.String()})
	switch results.Err() {
	case nil:
		// do nothing
	case mongo.ErrNoDocuments:
		return entity.AnalysisResults{}, repository.ErrAnalysisNoResults
	default:
		log.WithError(results.Err()).Errorf("error searching analysis for project_id: %v", projectID)
		return entity.AnalysisResults{}, repository.ErrAnalysisUnexpected
	}

	var dto analysisDTO
	if err := results.Decode(&dto); err != nil {
		log.WithError(err).Errorf("error decoding results for analysis with project_id: %v", projectID)
		return entity.AnalysisResults{}, repository.ErrAnalysisUnexpected
	}

	return adb.mapper.toEntity(dto), nil
}

// Delete removes an existing Analysis from the underlying MongoDB collection.
func (adb *AnalysisDB) Delete(ctx context.Context, id uuid.UUID) error {
	results, err := adb.collection.DeleteOne(ctx, bson.M{"_id": id.String()})
	if err != nil {
		log.WithError(err).Errorf("error deleting analysis with id: %v", id)
		return repository.ErrAnalysisUnexpected
	}

	if results.DeletedCount == 0 {
		return repository.ErrAnalysisNoResults
	}

	return nil
}
