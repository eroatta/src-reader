package extractors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: use a real test... not this piece of shit
func TestProcess_OnSamuraiExtractor_ShouldNotFail(t *testing.T) {
	samurai := NewSamuraiExtractor()

	Process(samurai, nil)

	assert.True(t, true)
}

func TestNewSamurai_ShouldReturnNewExtractor(t *testing.T) {
	extractor := NewSamuraiExtractor()

	assert.NotNil(t, extractor)
	assert.IsType(t, SamuraiExtractor{}, extractor)
}

func TestGetName_OnSamurai_ShouldReturnSamurai(t *testing.T) {
	extractor := NewSamuraiExtractor()

	assert.Equal(t, "samurai", extractor.Name())
}

func TestVisit_OnSamuraiWithNilNode_ShouldReturnNil(t *testing.T) {
	extractor := NewSamuraiExtractor()

	assert.Nil(t, extractor.Visit(nil))
}

func TestVisit_OnSamuraiWithNotNilNode_ShouldReturnVisitor(t *testing.T) {
	extractor := NewSamuraiExtractor()

	node, _ := parser.ParseExpr("a + b")
	got := extractor.Visit(node)

	assert.NotNil(t, got)
}

func TestVisit_OnSamurai_ShouldSplitTheIdentifiers(t *testing.T) {
	var tests = []struct {
		name        string
		src         string
		uniqueWords int
		expected    map[string]int
	}{
		{
			name: "VarDeclNode",
			src: `
				package main
				var testIden string
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"test": 1,
				"iden": 1,
			},
		},
		{
			name: "VarDelNode_MultipleNames",
			src: `
				package main
				var testIden, testBar string
			`,
			uniqueWords: 3,
			expected: map[string]int{
				"test": 2,
				"iden": 1,
				"bar":  1,
			},
		},
		{
			name: "VarDelNode_BlockAndValues",
			src: `
				package main
				var (
					testIden = "boo"
					foo = "bar"
				)
			`,
			uniqueWords: 5,
			expected: map[string]int{
				"test": 1,
				"iden": 1,
				"boo":  1,
				"foo":  1,
				"bar":  1,
			},
		},
		{
			name: "ConstDeclNode",
			src: `
				package main
				const testIden = "boo"
			`,
			uniqueWords: 3,
			expected: map[string]int{
				"test": 1,
				"iden": 1,
				"boo":  1,
			},
		},
		{
			name: "ConstDeclNode_MultipleNamesAndValues",
			src: `
				package main
				const testIden, foo = "boo", "bar"
			`,
			uniqueWords: 5,
			expected: map[string]int{
				"test": 1,
				"iden": 1,
				"boo":  1,
				"foo":  1,
				"bar":  1,
			},
		},
		{
			name: "ConstDeclNode_BlockAndValues",
			src: `
				package main
				const (
					testIden = "boo"
					foo = "bar"
				)
			`,
			uniqueWords: 5,
			expected: map[string]int{
				"test": 1,
				"iden": 1,
				"boo":  1,
				"foo":  1,
				"bar":  1,
			},
		},
		{
			name: "FuncDeclNode",
			src: `
				package main
				
				import "fmt"

				func main() {
					fmt.Println("Hello Samurai Extractor!")
				}
			`,
			uniqueWords: 1,
			expected: map[string]int{
				"main": 1,
			},
		},
		{
			name: "FuncDeclNode_Parameters",
			src: `
				package main

				import (
					"fmt"
					"gin"
				)
		
				func main() {
					engine := gin.New()
					engine.Delims("first_argument", "second_argument")
				}
		
				func (engine *Engine) Delims(left, right string) *Engine {
					engine.delims = render.Delims{Left: left, Right: right}
					return engine
				}
			`,
			uniqueWords: 4,
			expected: map[string]int{
				"main":   1,
				"delims": 1,
				"left":   1,
				"right":  1,
			},
		},
		{
			name: "FuncDeclNode_NamedResults",
			src: `
				package main

				import (
					"fmt"
				)
				
				func main() {
					text, number := results()
					fmt.Println(fmt.Sprintf("%s, %d", text, number))
				}
				
				func results() (text string, number int) {
					return "text", 10
				}
			`,
			uniqueWords: 4,
			expected: map[string]int{
				"main":    1,
				"results": 1,
				"text":    1,
				"number":  1,
			},
		},
		{
			name: "TypeDeclNode",
			src: `
				package main

				type myInt int
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"my":  1,
				"int": 1,
			},
		},
		{
			name: "StructTypeDeclNode_NoFields",
			src: `
				package main

				type myStruct struct {}
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"my":     1,
				"struct": 1,
			},
		},
		{
			name: "StructTypeDeclNode_TwoFields",
			src: `
				package main

				type myStruct struct {
					first string,
					second string
				}
			`,
			uniqueWords: 4,
			expected: map[string]int{
				"my":     1,
				"struct": 1,
				"first":  1,
				"second": 1,
			},
		},
		{
			name: "InterfaceTypeDeclNode",
			src: `
				package main

				type myInterface interface
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"my":        1,
				"interface": 1,
			},
		},
		{
			name: "InterfaceTypeDeclNode_TwoMethods",
			src: `
				package main

				type myInterface interface {
					myMethod(arg string) (out string)
					anotherMethod(arg string) (out string)
				}
			`,
			uniqueWords: 6,
			expected: map[string]int{
				"my":        2,
				"interface": 1,
				"method":    2,
				"another":   1,
				"arg":       2,
				"out":       2,
			},
		},
		{
			name: "FileCommentsNode",
			src: `
				// package comment
				package main
			
				// regular comment
				type MyInterface interface {
					MyMethod(out string) // inline comment
				}
			
				/* 
				block comment
				*/
				type abc int
			`,
			uniqueWords: 10,
			expected: map[string]int{
				"my":        2,
				"interface": 1,
				"method":    1,
				"out":       1,
				"abc":       1,
				"package":   1,
				"regular":   1,
				"inline":    1,
				"block":     1,
				"comment":   4,
			},
		},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			samurai := NewSamuraiExtractor()

			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			ast.Walk(samurai, node)

			assert.NotNil(t, samurai)

			extractor := samurai.(SamuraiExtractor)

			assert.NotEmpty(t, extractor.words)
			assert.Equal(t, fixture.uniqueWords, len(extractor.words))
			for key, value := range fixture.expected {
				assert.Equal(t, value, extractor.words[key], fmt.Sprintf("invalid number of occurrencies for element: %s", key))
			}
		})
	}
}
