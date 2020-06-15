package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
func (idb *IdentifierDB) Add(ctx context.Context, analysis entity.AnalysisResults, ident entity.Identifier) error {
	dto := idb.mapper.toDTO(ident, analysis)
	_, err := idb.collection.InsertOne(ctx, dto)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error inserting record %v", ident))
		return repository.ErrIdentifierUnexpected
	}

	return nil
}

// FindAllByProject retrieves all the identifiers related to a given project, from the underlying MongoDB collection.
func (idb *IdentifierDB) FindAllByProject(ctx context.Context, projectRef string) ([]entity.Identifier, error) {
	// TODO: review how to handle URL
	cursor, err := idb.collection.Find(ctx, bson.M{"project_ref": strings.TrimPrefix(projectRef, "https://github.com/")})
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error looking documents for %s", projectRef))
		return []entity.Identifier{}, repository.ErrIdentifierUnexpected
	}

	var elements []identifierDTO
	err = cursor.All(ctx, &elements)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error decoding found documents for %s", projectRef))
		return []entity.Identifier{}, repository.ErrIdentifierUnexpected
	}

	identifiers := make([]entity.Identifier, len(elements))
	for i, element := range elements {
		identifiers[i] = idb.mapper.toEntity(element)
	}

	return identifiers, nil
}

// FindAllByProject retrieves all the identifiers related to a given project, from the underlying MongoDB collection.
func (idb *IdentifierDB) FindAllByProjectAndFile(ctx context.Context, projectRef string, filename string) ([]entity.Identifier, error) {
	// TODO: review how to handle URL
	cursor, err := idb.collection.Find(ctx,
		bson.M{"project_ref": strings.TrimPrefix(projectRef, "https://github.com/"), "file": filename})
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error looking documents for %s on file %s", projectRef, filename))
		return []entity.Identifier{}, repository.ErrIdentifierUnexpected
	}

	var elements []identifierDTO
	err = cursor.All(ctx, &elements)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error decoding found documents for %s on file %s", projectRef, filename))
		return []entity.Identifier{}, repository.ErrIdentifierUnexpected
	}

	identifiers := make([]entity.Identifier, len(elements))
	for i, element := range elements {
		identifiers[i] = idb.mapper.toEntity(element)
	}

	return identifiers, nil
}
