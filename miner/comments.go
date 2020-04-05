package miner

import (
	"go/ast"
	"strings"

	"github.com/eroatta/src-reader/entity"
)

// Comments handles the comments mining process.
type Comments struct {
	comments []string
}

// NewComments creates a new Comments miner.
func NewComments() *Comments {
	return &Comments{
		comments: make([]string, 0),
	}
}

// Type returns the miner type.
func (c *Comments) Type() entity.MinerType {
	return entity.MinerComments
}

// Visit implements the ast.Visitor interface and handles the logic for the data extraction.
func (c *Comments) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	commentGroup, ok := node.(*ast.CommentGroup)
	if ok {
		for _, comment := range commentGroup.List {
			cleanComment := strings.Trim(cleaner.ReplaceAllString(comment.Text, " "), " ")
			if cleanComment == "" {
				continue
			}

			c.comments = append(c.comments, cleanComment)
		}
	}

	return c
}

// Collected returns the list of comments found.
func (c *Comments) Collected() []string {
	return c.comments
}
