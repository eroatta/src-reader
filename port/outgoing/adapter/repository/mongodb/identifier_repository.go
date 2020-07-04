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

// FindAllByAnalysisID retrieves all the identifiers related to a given analysis, from the underlying MongoDB collection.
func (idb *IdentifierDB) FindAllByAnalysisID(ctx context.Context, analysisID uuid.UUID) ([]entity.Identifier, error) {
	cursor, err := idb.collection.Find(ctx, bson.M{"analysis_id": analysisID.String()})
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error looking documents for analysis ID %v", analysisID))
		return []entity.Identifier{}, repository.ErrIdentifierUnexpected
	}

	var elements []identifierDTO
	err = cursor.All(ctx, &elements)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("error decoding found documents for analysis ID %v", analysisID))
		return []entity.Identifier{}, repository.ErrIdentifierUnexpected
	}

	identifiers := make([]entity.Identifier, len(elements))
	for i, element := range elements {
		identifiers[i] = idb.mapper.toEntity(element)
	}

	return identifiers, nil
}

// FindAllByProjectAndFile retrieves all the identifiers for a file related to a given project, from the underlying MongoDB collection.
func (idb *IdentifierDB) FindAllByProjectAndFile(ctx context.Context, projectRef string, filename string) ([]entity.Identifier, error) {
	cursor, err := idb.collection.Find(ctx, bson.M{"project_ref": projectRef, "file": filename})
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

// DeleteAllByAnalysisID removes a set of existing identifires from the underlying MongoDB collection.
func (idb *IdentifierDB) DeleteAllByAnalysisID(ctx context.Context, analysisID uuid.UUID) error {
	results, err := idb.collection.DeleteMany(ctx, bson.M{"analysis_id": analysisID.String()})
	if err != nil {
		log.WithError(err).Errorf("error deleting identifiers with analysis_id: %v", analysisID)
		return repository.ErrIdentifierUnexpected
	}

	if results.DeletedCount == 0 {
		return repository.ErrIdentifierNoResults
	}

	return nil
}
