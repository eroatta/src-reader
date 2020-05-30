package mongodb

import (
	"context"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const insightCollection string = "insight"

// InsightDB represents a MongoDB database, focused on the collection handling the insights documents.
type InsightDB struct {
	client    *mongo.Client
	mapper    *insightMapper
	collecton *mongo.Collection
}

// NewMongoDBInsightRepository creates a repository.InsightRepository backed up by a MongoDB database.
func NewMongoDBInsightRepository(client *mongo.Client, dbname string) *InsightDB {
	return &InsightDB{
		client:    client,
		mapper:    &insightMapper{},
		collecton: client.Database(dbname).Collection(insightCollection),
	}
}

// AddAll transforms and stores a set of entity.Insight entities into documents on the underlying
// MongoDB collection.
func (idb *InsightDB) AddAll(ctx context.Context, insights []entity.Insight) error {
	dtos := make([]interface{}, 0)
	for _, insight := range insights {
		dtos = append(dtos, idb.mapper.toDTO(insight))
	}

	results, err := idb.collecton.InsertMany(ctx, dtos)
	if err != nil {
		log.WithError(err).Error("error inserting records")
		return repository.ErrInsightUnexpected
	}
	log.WithField("count", len(results.InsertedIDs)).Debug("inserted insights")

	return nil
}
