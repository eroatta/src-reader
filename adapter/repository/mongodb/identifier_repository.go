package mongodb

import (
	"context"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const identifiersCollection string = "identifiers"

// IdentifierDB represents a MongoDB database, focused on the collection handling the identifiers documents.
type IdentifierDB struct {
	client     *mongo.Client
	mapper     *identifierMapper
	collection *mongo.Collection
}

// NewMongoDBIdentifierRepository creates a repository.IdentifierRepository backed up by a MongoDB database.
func NewMongoDBIdentifierRepository(client *mongo.Client, dbname string) *IdentifierDB {
	return &IdentifierDB{
		client:     client,
		mapper:     &identifierMapper{},
		collection: client.Database(dbname).Collection(identifiersCollection),
	}
}

// Add transforms and stores an Identifier entity into a document on the underlying MongoDB collection.
func (idb *IdentifierDB) Add(ctx context.Context, project entity.Project, ident entity.Identifier) error {
	dto := idb.mapper.toDTO(ident, project)
	_, err := idb.collection.InsertOne(ctx, dto)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error inserting record %v", ident))
		return repository.ErrIdentifierUnexpected
	}

	return nil
}
