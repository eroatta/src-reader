package mongodb

import (
	"context"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
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
		return repository.ErrProjectUnexpected // TODO: change for specific repository error
	}

	return nil
}
