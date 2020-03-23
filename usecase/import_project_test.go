package usecase_test

import (
	"testing"

	"github.com/eroatta/src-reader/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewImportProjectUsecase_ShouldReturnNewInstance(t *testing.T) {
	uc := usecase.NewImportProjectUsecase()

	assert.NotNil(t, uc)
}

func TestImport_OnImportProjectUsecase_ShouldReturnImportResults(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestImport_OnImportProjectUsecase_WhenAlreadyImportedProject_ShouldImportResults(t *testing.T) {
	// TODO: check if should fail or retrieve previous results...
	assert.FailNow(t, "not yet implemented")
}

func TestImport_OnImportProjectUsecase_WhenUnableToCheckExistingProject_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestImport_OnImportProjectUsecase_WhenUnableToRetrieveMetadataFromRemoteRepository_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestImport_OnImportProjectUsecase_WhenUnableToCloneSourceCode_ShouldReturnError(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}
