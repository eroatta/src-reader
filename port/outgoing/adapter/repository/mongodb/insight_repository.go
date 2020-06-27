package mongodb

import (
	"context"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const insightCollection string = "insight"

// InsightDB represents a MongoDB database, focused on the collection handling the insights documents.
type InsightDB struct {
	client     *mongo.Client
	mapper     *insightMapper
	collection *mongo.Collection
}

// NewMongoDBInsightRepository creates a repository.InsightRepository backed up by a MongoDB database.
func NewMongoDBInsightRepository(client *mongo.Client, dbname string) *InsightDB {
	return &InsightDB{
		client:     client,
		mapper:     &insightMapper{},
		collection: client.Database(dbname).Collection(insightCollection),
	}
}

// AddAll transforms and stores a set of entity.Insight entities into documents on the underlying
// MongoDB collection.
func (idb *InsightDB) AddAll(ctx context.Context, insights []entity.Insight) error {
	dtos := make([]interface{}, 0)
	for _, insight := range insights {
		dtos = append(dtos, idb.mapper.toDTO(insight))
	}

	results, err := idb.collection.InsertMany(ctx, dtos)
	if err != nil {
		log.WithError(err).Error("error inserting records")
		return repository.ErrInsightUnexpected
	}
	log.WithField("count", len(results.InsertedIDs)).Debug("inserted insights")

	return nil
}

// GetByAnalysisID finds a set of existing insights on the underlying MongoDB collection, and returns
// them as entity.Insight.
func (idb *InsightDB) GetByAnalysisID(ctx context.Context, analysisID uuid.UUID) ([]entity.Insight, error) {
	cursor, err := idb.collection.Find(ctx, bson.M{"analysis_id": analysisID.String()})
	switch err {
	case nil:
		// do nothing
	case mongo.ErrNoDocuments:
		return []entity.Insight{}, repository.ErrInsightNoResults
	default:
		log.WithError(err).Errorf("error searching insights with analysis_id: %v", analysisID)
		return []entity.Insight{}, repository.ErrInsightUnexpected
	}

	var elements []insightDTO
	err = cursor.All(ctx, &elements)
	if err != nil {
		log.WithError(err).Errorf("error decoding found documents for analysis_id: %v", analysisID)
		return []entity.Insight{}, repository.ErrInsightUnexpected
	}

	insights := make([]entity.Insight, len(elements))
	for i, element := range elements {
		insights[i] = idb.mapper.toEntity(element)
	}

	return insights, nil
}

// DeleteAllByAnalysisID removes a set of existing insights from the underlying MongoDB collection.
func (idb *InsightDB) DeleteAllByAnalysisID(ctx context.Context, analysisID uuid.UUID) error {
	results, err := idb.collection.DeleteMany(ctx, bson.M{"analysis_id": analysisID.String()})
	if err != nil {
		log.WithError(err).Errorf("error deleting insights with analysis_id: %v", analysisID)
		return repository.ErrInsightUnexpected
	}

	if results.DeletedCount == 0 {
		return repository.ErrInsightNoResults
	}

	return nil
}
