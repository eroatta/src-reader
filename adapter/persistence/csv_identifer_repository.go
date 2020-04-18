package persistence

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/eroatta/src-reader/entity"
	log "github.com/sirupsen/logrus"
)

// NewCSVIdentifierRepository creates a new CSVIdentifierRepository to store
// identifiers as comma-separated values on the referenced file.
func NewCSVIdentifierRepository(file *os.File) *CSVIdentifierRepository {
	return &CSVIdentifierRepository{
		output: csv.NewWriter(file),
	}
}

// CSVIdentifierRepository represents a repository capable of storing identifiers
// as comma-separated values.
type CSVIdentifierRepository struct {
	output *csv.Writer
}

// Add creates a new row for an identifier, using comma-separated values. The columns added are:
// * id
// * file       string
//	Position   token.Pos
//	Name       string
//	Type       token.Token
//	Parent     string
//	ParentPos  token.Pos
//	Node       *ast.Node
//	Splits     map[string][]Split
//	Expansions map[string][]Expansion
//	Error      error
func (r *CSVIdentifierRepository) Add(ctx context.Context, project entity.Project, ident entity.Identifier) error {
	var splitsColumnBuilder strings.Builder
	for sptr, splits := range ident.Splits {
		softwords := make([]string, len(splits))
		for _, split := range splits {
			softwords[split.Order-1] = split.Value
		}

		splitsColumnBuilder.WriteString(sptr)
		splitsColumnBuilder.WriteString(":")
		splitsColumnBuilder.WriteString(strings.Join(softwords, "-"))
		splitsColumnBuilder.WriteString("|")
	}
	splits := strings.TrimSuffix(splitsColumnBuilder.String(), "|")

	var expansionsColumnBuilder strings.Builder
	for expr, expansions := range ident.Expansions {
		softwords := make([]string, len(expansions))
		for i, expansion := range expansions {
			softwords[i] = fmt.Sprintf("%s = %s", expansion.From, strings.Join(expansion.Values, "_"))
		}

		expansionsColumnBuilder.WriteString(expr)
		expansionsColumnBuilder.WriteString(":")
		expansionsColumnBuilder.WriteString(strings.Join(softwords, "+++"))
		expansionsColumnBuilder.WriteString("|")
	}
	expansions := strings.TrimSuffix(expansionsColumnBuilder.String(), "|")

	row := []string{
		ident.ID,
		ident.File,
		fmt.Sprintf("%v", ident.Position),
		ident.Type.String(),
		fmt.Sprintf("%v", ident.IsLocal()),
		ident.Parent,
		fmt.Sprintf("%v", ident.ParentPos),
		splits,
		expansions,
		printableError(ident.Error),
	}

	if err := r.output.Write(row); err != nil {
		log.WithError(err).WithField("row", row).Error(fmt.Sprintf("unable to store identifier with ID %s", ident.ID))
		return err
	}

	r.output.Flush()
	return nil
}

func printableError(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}
