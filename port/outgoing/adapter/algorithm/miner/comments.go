package miner

import (
	"go/ast"
	"strings"

	"github.com/eroatta/src-reader/entity"
)

// NewCommentsFactory creates a new comments miner factory.
func NewCommentsFactory() entity.MinerFactory {
	return commentsFactory{}
}

type commentsFactory struct{}

func (f commentsFactory) Make() (entity.Miner, error) {
	return NewComments(), nil
}

// NewComments creates a new Comments miner.
func NewComments() *Comments {
	return &Comments{
		miner:    miner{"comments"},
		comments: make([]string, 0),
	}
}

// Comments handles the comments mining process.
type Comments struct {
	miner
	comments []string
}

// SetCurrentFile specifies the current file being mined.
func (c *Comments) SetCurrentFile(filename string) {
	// do nothing
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

// Results returns the list of comments found.
func (c *Comments) Results() interface{} {
	return c.comments
}
