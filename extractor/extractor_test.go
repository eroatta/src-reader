package extractor_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/code"
	"github.com/eroatta/src-reader/extractor"
	"github.com/stretchr/testify/assert"
)

func TestVisit_OnExtractorWithFuncDecl_ShouldReturnFoundIdentifiers(t *testing.T) {
	src := `
		package main

		func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
			path += root.path
			if len(root.handlers) > 0 {
				handlerFunc := root.handlers.Last()
				routes = append(routes, RouteInfo{
					Method:      method,
					Path:        path,
					Handler:     nameOfFunction(handlerFunc),
					HandlerFunc: handlerFunc,
				})
			}
			for _, child := range root.children {
				routes = iterate(path, method, routes, child)
			}
			return routes
		}
	`

	expected := []code.Identifier{
		{
			File:       "testfile",
			Position:   20,
			Name:       "iterate",
			Type:       "FuncDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   154,
			Name:       "handlerFunc",
			Type:       "AssignStmt",
			Parent:     "iterate",
			ParentPos:  20,
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	e := extractor.New("testfile")
	ast.Walk(e, node.Decls[0])

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}

func TestVisit_OnExtractorWithFuncDeclUsingSameFuncName_ShouldReturnFoundIdentifiers(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestVisit_OnExtractorWithFuncDecl_ShouldReturnFoundLocalIdentifiers(t *testing.T) {
	assert.FailNow(t, "not yet implemented")
}

func TestVisit_OnExtractorWithVarDecl_ShouldReturnFoundIdentifiers(t *testing.T) {
	src := `
		package main
		
		var (
			common string
			regular string = "valid"
			nrzXXZ int = 32
		)
	`

	expected := []code.Identifier{
		{
			File:       "testfile",
			Position:   31,
			Name:       "common",
			Type:       "VarDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   48,
			Name:       "regular",
			Type:       "VarDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   76,
			Name:       "nrzXXZ",
			Type:       "VarDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	e := extractor.New("testfile")
	ast.Walk(e, node.Decls[0])

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}

func TestVisit_OnExtractorWithConstDecl_ShouldReturnFoundIdentifiers(t *testing.T) {
	src := `
		package main
		
		// outer comment
		const (
			common string = "common"
			regular, notRegular string = "valid", "invalid"
			nrzXXZ int = 32
		)
	`

	expected := []code.Identifier{
		{
			File:       "testfile",
			Position:   52,
			Name:       "common",
			Type:       "ConstDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   80,
			Name:       "regular",
			Type:       "ConstDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   89,
			Name:       "notRegular",
			Type:       "ConstDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   131,
			Name:       "nrzXXZ",
			Type:       "ConstDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	e := extractor.New("testfile")
	ast.Walk(e, node.Decls[0])

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}

func TestVisit_OnExtractorWithStructDecl_ShouldReturnFoundIdentifiers(t *testing.T) {
	src := `
		package main
		
		type (
			// local comment
			selector struct {
				pick string
			}
		
			httpClient struct {
				protocolPicker string
				url string
			}
		)
	`

	expected := []code.Identifier{
		{
			File:       "testfile",
			Position:   52,
			Name:       "selector",
			Type:       "StructDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   97,
			Name:       "httpClient",
			Type:       "StructDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	e := extractor.New("testfile")
	ast.Walk(e, node.Decls[0])

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}

func TestVisit_OnExtractorWithInterfaceDecl_ShouldReturnFoundIdentifiers(t *testing.T) {
	src := `
		package main
		
		type (
			// local comment
			selector interface {
				pick() string
			}
		
			httpClient interface {
				protocolPicker() string
				url() string
			}
		)
	`

	expected := []code.Identifier{
		{
			File:       "testfile",
			Position:   52,
			Name:       "selector",
			Type:       "InterfaceDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
		{
			File:       "testfile",
			Position:   102,
			Name:       "httpClient",
			Type:       "InterfaceDecl",
			Splits:     make(map[string][]string),
			Expansions: make(map[string][]string),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile", []byte(src), parser.ParseComments)

	e := extractor.New("testfile")
	ast.Walk(e, node.Decls[0])

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}
