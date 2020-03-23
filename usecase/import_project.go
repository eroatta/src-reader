package usecase

import (
	"context"
	"errors"

	"github.com/eroatta/src-reader/repository"
)

var (
	ErrUnableToReadProject      = errors.New("")
	ErrUnableToRetrieveMetadata = errors.New("")
	ErrUnableToCloneSourceCode  = errors.New("")
	ErrUnableToSaveProject      = errors.New("")
)

// ImportResults represents the results after the project import process.
type ImportResults struct {
	ProjectURL  string
	Status      string
	ErrorDetail error
}

// ImportProjectUsecase defines the contract for the use case related to the import process
// for a project.
type ImportProjectUsecase interface {
	// Import executes the pipeline to import a project from GitHub and
	// extract the identifiers. It returns the process results.
	Import(ctx context.Context, url string) (ImportResults, error)
}

// NewImportProjectUsecase initializes a new ImportProjectUsecase handler.
func NewImportProjectUsecase() ImportProjectUsecase {
	return importProjectUsecase{}
}

type importProjectUsecase struct {
	projectRepository    repository.ProjectRepository
	remoteRepository     repository.RemoteProjectRepository
	sourceCodeRepository repository.SourceCodeRepository
}

func (uc importProjectUsecase) Import(ctx context.Context, url string) (ImportResults, error) {
	// check if not previously imported
	project, err := uc.projectRepository.GetByURL(ctx, url)
	if err != nil {
		// TODO: handle error
	}

	// TODO: nil or empty
	if project != (repository.Project{}) {
		// TODO: create ImportResults
	}

	// retrieve metadata
	metadata, err := uc.remoteRepository.RetrieveMetadata(ctx, url)
	if err != nil {
		// TODO: handle error
	}

	project = repository.Project{
		Status:   "retrieved",
		Metadata: metadata,
	}

	// clone the source code
	err = uc.sourceCodeRepository.Clone(ctx, url)
	if err != nil {
		// TODO: handle error
	}

	// store the results
	err = uc.projectRepository.Add(ctx, project)
	if err != nil {
		// TODO: handle error
	}

	return ImportResults{}, nil
}
