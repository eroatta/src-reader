package extractor_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/eroatta/src-reader/entity"
	"github.com/eroatta/src-reader/port/outgoing/adapter/algorithm/extractor"
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

	expected := []entity.Identifier{
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:func+++name:iterate",
			Package:    "main",
			File:       "testfile.go",
			Position:   20,
			Name:       "iterate",
			Type:       token.FUNC,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	e := extractor.New("testfile.go")
	ast.Walk(e, node)

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}

func TestVisit_OnExtractorWithFuncDeclUsingSameFuncName_ShouldReturnFoundIdentifiers(t *testing.T) {
	src := `
		package main

		type car struct {}

		func (c car) name() {
			// do nothing
		}

		type boat struct{}

		func (b *boat) name() {
			// do nothing
		}
	`

	expected := []entity.Identifier{
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:struct+++name:car",
			Package:    "main",
			File:       "testfile.go",
			Position:   25,
			Name:       "car",
			Type:       token.STRUCT,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:func+++name:car.name",
			Package:    "main",
			File:       "testfile.go",
			Position:   42,
			Name:       "name",
			Type:       token.FUNC,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:struct+++name:boat",
			Package:    "main",
			File:       "testfile.go",
			Position:   93,
			Name:       "boat",
			Type:       token.STRUCT,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:func+++name:boat.name",
			Package:    "main",
			File:       "testfile.go",
			Position:   110,
			Name:       "name",
			Type:       token.FUNC,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	e := extractor.New("testfile.go")
	ast.Walk(e, node)

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
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

	expected := []entity.Identifier{
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:var+++name:common",
			Package:    "main",
			File:       "testfile.go",
			Position:   31,
			Name:       "common",
			Type:       token.VAR,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:var+++name:regular",
			Package:    "main",
			File:       "testfile.go",
			Position:   48,
			Name:       "regular",
			Type:       token.VAR,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:var+++name:nrzXXZ",
			Package:    "main",
			File:       "testfile.go",
			Position:   76,
			Name:       "nrzXXZ",
			Type:       token.VAR,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	e := extractor.New("testfile.go")
	ast.Walk(e, node)

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

	expected := []entity.Identifier{
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:const+++name:common",
			Package:    "main",
			File:       "testfile.go",
			Position:   52,
			Name:       "common",
			Type:       token.CONST,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:const+++name:regular",
			Package:    "main",
			File:       "testfile.go",
			Position:   80,
			Name:       "regular",
			Type:       token.CONST,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:const+++name:notRegular",
			Package:    "main",
			File:       "testfile.go",
			Position:   89,
			Name:       "notRegular",
			Type:       token.CONST,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:const+++name:nrzXXZ",
			Package:    "main",
			File:       "testfile.go",
			Position:   131,
			Name:       "nrzXXZ",
			Type:       token.CONST,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	e := extractor.New("testfile.go")
	ast.Walk(e, node)

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

	expected := []entity.Identifier{
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:struct+++name:selector",
			Package:    "main",
			File:       "testfile.go",
			Position:   52,
			Name:       "selector",
			Type:       token.STRUCT,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:struct+++name:httpClient",
			Package:    "main",
			File:       "testfile.go",
			Position:   97,
			Name:       "httpClient",
			Type:       token.STRUCT,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	e := extractor.New("testfile.go")
	ast.Walk(e, node)

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

	expected := []entity.Identifier{
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:interface+++name:selector",
			Package:    "main",
			File:       "testfile.go",
			Position:   52,
			Name:       "selector",
			Type:       token.INTERFACE,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
		{
			ID:         "filename:testfile.go+++pkg:main+++declType:interface+++name:httpClient",
			Package:    "main",
			File:       "testfile.go",
			Position:   102,
			Name:       "httpClient",
			Type:       token.INTERFACE,
			Splits:     make(map[string][]entity.Split),
			Expansions: make(map[string][]entity.Expansion),
		},
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "testfile.go", []byte(src), parser.ParseComments)

	e := extractor.New("testfile.go")
	ast.Walk(e, node)

	identifiers := e.Identifiers()
	assert.Equal(t, expected, identifiers)
}
