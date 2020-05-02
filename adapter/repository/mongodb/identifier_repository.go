package mongodb

import (
	"context"
	"fmt"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdentifierDB struct {
	client     *mongo.Client
	mapper     *identifierMapper
	collection *mongo.Collection
}

func NewMongoDBIdentifierRepository(client *mongo.Client, dbname string) *IdentifierDB {
	return &IdentifierDB{
		client:     client,
		mapper:     &identifierMapper{},
		collection: client.Database(dbname).Collection("identifiers"),
	}
}

func (idb *IdentifierDB) Add(ctx context.Context, project entity.Project, ident entity.Identifier) error {
	dto := idb.mapper.toDTO(ident, project)
	_, err := idb.collection.InsertOne(ctx, dto)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error inserting record %v", ident))
		return repository.ErrIdentifierUnexpected
	}

	return nil
}
