package usecase

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ast/astutil"
)

type RewrittenFileUsecase interface {
	Process(ctx context.Context, projectRef string, filename string) ([]byte, error)
}

func NewRewrittenFileUsecase(pr repository.ProjectRepository, scr repository.SourceCodeRepository, ir repository.IdentifierRepository) RewrittenFileUsecase {
	return &rewrittenFileUsecase{
		pr:  pr,
		scr: scr,
		ir:  ir,
	}
}

type rewrittenFileUsecase struct {
	pr  repository.ProjectRepository
	scr repository.SourceCodeRepository
	ir  repository.IdentifierRepository
}

func (uc *rewrittenFileUsecase) Process(ctx context.Context, projectRef string, filename string) ([]byte, error) {
	project, err := uc.pr.GetByReference(ctx, projectRef)
	switch err {
	case nil:
		// do nothing
	case repository.ErrProjectNoResults:
		return nil, ErrProjectNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve project %s", projectRef)
		return nil, ErrUnexpected
	}

	found := false
	for _, f := range project.SourceCode.Files {
		if f == filename {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrFileNotFound
	}

	raw, err := uc.scr.Read(ctx, project.SourceCode.Location, filename)
	if err != nil {
		log.WithError(err).Errorf("unable to read file %s on project %s", filename, projectRef)
		return nil, ErrUnexpected
	}

	// retrieve identifiers
	// TODO: review how to handle URLs
	identifiers, err := uc.ir.FindAllByProjectAndFile(ctx, fmt.Sprintf("https://github.com/%s", projectRef), filename)
	switch err {
	case nil:
		// do nothing
	case repository.ErrIdentifierNoResults:
		return nil, ErrIdentifiersNotFound
	default:
		log.WithError(err).Errorf("unable to retrieve identifiers for project %s and file %s", projectRef, filename)
		return nil, ErrUnexpected
	}

	rename := make(map[string]entity.Identifier)
	for _, identifier := range identifiers {
		if identifier.Name != identifier.Normalization.Word {
			rename[identifier.Name] = identifier
		}
	}

	// parse file
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, filename, raw, parser.ParseComments)
	if err != nil {
		log.WithError(err).Errorf("unable to parse file %s from project %s", filename, projectRef)
		return nil, ErrUnexpected
	}

	// make a copy and rewrite AST with identifiers
	rewritten := astutil.Apply(f, func(cr *astutil.Cursor) bool {
		ident, ok := cr.Node().(*ast.Ident)
		if !ok {
			return true
		}

		newIdent, ok := rename[ident.Name]
		if !ok {
			return true
		}

		if obj, ok := f.Scope.Objects[ident.Name]; ok && ident.Obj == obj {
			ident.Name = newIdent.Normalization.Word
		}

		return true
	}, nil)

	var modif bytes.Buffer
	if err = format.Node(&modif, fs, rewritten); err != nil {
		log.WithError(err).Errorf("unable to rewritte file %s from project %s", filename, projectRef)
		return nil, ErrUnexpected
	}

	return modif.Bytes(), nil
}
