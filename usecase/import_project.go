package usecase

import "context"

// ImportResults represents the results after the project import process.
type ImportResults struct {
	Status string
}

type ImportProjectUseCase interface {
	// Import executes the pipeline to import a project from GitHub and
	// extract the identifiers. It returns the process results.
	Import(ctx context.Context, url string) (ImportResults, error)
}

func NewImportProjectUseCase() importProjectUseCase {
	return importProjectUseCase{}
}

type importProjectUseCase struct {
}

func (uc importProjectUseCase) Import(ctx context.Context, url string) (ImportResults, error) {
	return ImportResults{}, nil
}
