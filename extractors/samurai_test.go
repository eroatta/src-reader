package extractors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				
				func results() (text string, number int) {
					return "text", 10
				}
			`,
			uniqueWords: 3,
			expected: map[string]int{
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
			
				// tricky comment has code: var := obj.Call("/route/v1", "text")

				/* 
				block comment
				*/
				type abc int
			`,
			uniqueWords: 20,
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
				"tricky":    1,
				"has":       1,
				"code":      1,
				"var":       1,
				"obj":       1,
				"call":      1,
				"route":     1,
				"v":         1,
				"1":         1,
				"text":      1,
				"comment":   5,
			},
		},
		{
			name: "AssignStmt",
			src: `
				package main
			
				func main() {
					text := "define text"					
					text = "assign text"
				}
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"main": 1,
				"text": 1,
			},
		},
		{
			name: "AssignStmt_IgnoredVar",
			src: `
				package main
			
				func main() {
					_, err := check()
				}
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"main": 1,
				"err":  1,
			},
		},
		{
			name: "AssignStmt_ReusedVar",
			src: `
				package main
			
				import "errors"

				func main() {
					var err errors.error
					_, err := check()
				}
			`,
			uniqueWords: 2,
			expected: map[string]int{
				"main": 1,
				"err":  1,
			},
		},
		{
			name: "RangeStmt",
			src: `
				package main
			
				import "fmt"

				func main() {
					for index, value := range []string{"test"} {
						fmt.Sprintf("%d, %s", index, value)
					}
				}
			`,
			uniqueWords: 3,
			expected: map[string]int{
				"main":  1,
				"index": 1,
				"value": 1,
			},
		},
	}

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, _ := parser.ParseFile(fs, "", []byte(fixture.src), parser.ParseComments)

			samurai := NewSamuraiExtractor()
			ast.Walk(samurai, node)

			assert.NotEmpty(t, samurai.words)
			assert.Equal(t, fixture.uniqueWords, len(samurai.words))
			for key, value := range fixture.expected {
				assert.Equal(t, value, samurai.words[key], fmt.Sprintf("invalid number of occurrencies for element: %s", key))
			}
		})
	}
}

func TestVisit_OnSamuraiWithFullFile_ShouldSplitCommentsAndIdentifiers(t *testing.T) {
	src := `
		// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
		// Use of this source code is governed by a MIT style
		// license that can be found in the LICENSE file.
		
		package gin
		
		import (
			"net/http"
			"path"
			"regexp"
			"strings"
		)
		
		// IRouter defines all router handle interface includes single and group router.
		type IRouter interface {
			IRoutes
			Group(string, ...HandlerFunc) *RouterGroup
		}
		
		// IRoutes defines all router handle interface.
		type IRoutes interface {
			Use(...HandlerFunc) IRoutes
		
			Handle(string, string, ...HandlerFunc) IRoutes
			Any(string, ...HandlerFunc) IRoutes
			GET(string, ...HandlerFunc) IRoutes
			POST(string, ...HandlerFunc) IRoutes
			DELETE(string, ...HandlerFunc) IRoutes
			PATCH(string, ...HandlerFunc) IRoutes
			PUT(string, ...HandlerFunc) IRoutes
			OPTIONS(string, ...HandlerFunc) IRoutes
			HEAD(string, ...HandlerFunc) IRoutes
		
			StaticFile(string, string) IRoutes
			Static(string, string) IRoutes
			StaticFS(string, http.FileSystem) IRoutes
		}
		
		// RouterGroup is used internally to configure router, a RouterGroup is associated with
		// a prefix and an array of handlers (middleware).
		type RouterGroup struct {
			Handlers HandlersChain
			basePath string
			engine   *Engine
			root     bool
		}
		
		var _ IRouter = &RouterGroup{}
		
		// Use adds middleware to the group, see example code in GitHub.
		func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
			group.Handlers = append(group.Handlers, middleware...)
			return group.returnObj()
		}
		
		// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
		// For example, all the routes that use a common middleware for authorization could be grouped.
		func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
			return &RouterGroup{
				Handlers: group.combineHandlers(handlers),
				basePath: group.calculateAbsolutePath(relativePath),
				engine:   group.engine,
			}
		}
		
		// BasePath returns the base path of router group.
		// For example, if v := router.Group("/rest/n/v1/api"), v.BasePath() is "/rest/n/v1/api".
		func (group *RouterGroup) BasePath() string {
			return group.basePath
		}
		
		func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
			absolutePath := group.calculateAbsolutePath(relativePath)
			handlers = group.combineHandlers(handlers)
			group.engine.addRoute(httpMethod, absolutePath, handlers)
			return group.returnObj()
		}
		
		// Handle registers a new request handle and middleware with the given path and method.
		// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
		// See the example code in GitHub.
		//
		// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
		// functions can be used.
		//
		// This function is intended for bulk loading and to allow the usage of less
		// frequently used, non-standardized or custom methods (e.g. for internal
		// communication with a proxy).
		func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
			if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
				panic("http method " + httpMethod + " is not valid")
			}
			return group.handle(httpMethod, relativePath, handlers)
		}
		
		// POST is a shortcut for router.Handle("POST", path, handle).
		func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("POST", relativePath, handlers)
		}
		
		// GET is a shortcut for router.Handle("GET", path, handle).
		func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("GET", relativePath, handlers)
		}
		
		// DELETE is a shortcut for router.Handle("DELETE", path, handle).
		func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("DELETE", relativePath, handlers)
		}
		
		// PATCH is a shortcut for router.Handle("PATCH", path, handle).
		func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("PATCH", relativePath, handlers)
		}
		
		// PUT is a shortcut for router.Handle("PUT", path, handle).
		func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("PUT", relativePath, handlers)
		}
		
		// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
		func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("OPTIONS", relativePath, handlers)
		}
		
		// HEAD is a shortcut for router.Handle("HEAD", path, handle).
		func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
			return group.handle("HEAD", relativePath, handlers)
		}
		
		// Any registers a route that matches all the HTTP methods.
		// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
		func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
			group.handle("GET", relativePath, handlers)
			group.handle("POST", relativePath, handlers)
			group.handle("PUT", relativePath, handlers)
			group.handle("PATCH", relativePath, handlers)
			group.handle("HEAD", relativePath, handlers)
			group.handle("OPTIONS", relativePath, handlers)
			group.handle("DELETE", relativePath, handlers)
			group.handle("CONNECT", relativePath, handlers)
			group.handle("TRACE", relativePath, handlers)
			return group.returnObj()
		}
		
		// StaticFile registers a single route in order to serve a single file of the local filesystem.
		// router.StaticFile("favicon.ico", "./resources/favicon.ico")
		func (group *RouterGroup) StaticFile(relativePath, filepath string) IRoutes {
			if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
				panic("URL parameters can not be used when serving a static file")
			}
			handler := func(c *Context) {
				c.File(filepath)
			}
			group.GET(relativePath, handler)
			group.HEAD(relativePath, handler)
			return group.returnObj()
		}
		
		// Static serves files from the given file system root.
		// Internally a http.FileServer is used, therefore http.NotFound is used instead
		// of the Router's NotFound handler.
		// To use the operating system's file system implementation,
		// use :
		//     router.Static("/static", "/var/www")
		func (group *RouterGroup) Static(relativePath, root string) IRoutes {
			return group.StaticFS(relativePath, Dir(root, false))
		}
		
		// StaticFS works just like Static() but a custom http.FileSystem can be used instead.
		// Gin by default user: gin.Dir()
		func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
			if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
				panic("URL parameters can not be used when serving a static folder")
			}
			handler := group.createStaticHandler(relativePath, fs)
			urlPattern := path.Join(relativePath, "/*filepath")
		
			// Register GET and HEAD handlers
			group.GET(urlPattern, handler)
			group.HEAD(urlPattern, handler)
			return group.returnObj()
		}
		
		func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
			absolutePath := group.calculateAbsolutePath(relativePath)
			fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
		
			return func(c *Context) {
				if _, nolisting := fs.(*onlyfilesFS); nolisting {
					c.Writer.WriteHeader(http.StatusNotFound)
				}
		
				file := c.Param("filepath")
				// Check if file exists and/or if we have permission to access it
				if _, err := fs.Open(file); err != nil {
					c.Writer.WriteHeader(http.StatusNotFound)
					c.handlers = group.engine.allNoRoute
					// Reset index
					c.index = -1
					return
				}
		
				fileServer.ServeHTTP(c.Writer, c.Request)
			}
		}
		
		func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
			finalSize := len(group.Handlers) + len(handlers)
			if finalSize >= int(abortIndex) {
				panic("too many handlers")
			}
			mergedHandlers := make(HandlersChain, finalSize)
			copy(mergedHandlers, group.Handlers)
			copy(mergedHandlers[len(group.Handlers):], handlers)
			return mergedHandlers
		}
		
		func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
			return joinPaths(group.basePath, relativePath)
		}
		
		func (group *RouterGroup) returnObj() IRoutes {
			if group.root {
				return group.engine
			}
			return group
		}
	`

	expectedWords := map[string]int{
		"1":              2,
		"2014":           1,
		"a":              19,
		"absolute":       3,
		"access":         1,
		"add":            1,
		"adds":           1,
		"all":            6,
		"allow":          1,
		"almeida":        1,
		"among":          1,
		"an":             1,
		"and":            9,
		"any":            3,
		"api":            2,
		"array":          1,
		"associated":     1,
		"authorization":  1,
		"base":           5,
		"be":             7,
		"bulk":           1,
		"but":            1,
		"by":             2,
		"calculate":      1,
		"can":            4,
		"check":          1,
		"code":           3,
		"combine":        1,
		"common":         2,
		"communication":  1,
		"configure":      1,
		"connect":        1,
		"copyright":      1,
		"could":          1,
		"create":         1,
		"creates":        1,
		"custom":         2,
		"default":        1,
		"defines":        2,
		"delete":         6,
		"different":      1,
		"dir":            1,
		"e":              1,
		"engine":         1,
		"err":            2,
		"example":        4,
		"exists":         1,
		"favicon":        2,
		"file":           13,
		"filepath":       1,
		"files":          1,
		"filesystem":     1,
		"final":          1,
		"for":            13,
		"found":          3,
		"frequently":     1,
		"from":           1,
		"fs":             5,
		"function":       1,
		"functions":      1,
		"g":              1,
		"get":            7,
		"gin":            2,
		"git":            2,
		"given":          2,
		"governed":       1,
		"group":          11,
		"grouped":        1,
		"handle":         21,
		"handler":        6,
		"handlers":       17,
		"have":           2,
		"head":           6,
		"http":           6,
		"hub":            2,
		"i":              4,
		"ico":            2,
		"if":             3,
		"implementation": 1,
		"in":             4,
		"includes":       1,
		"index":          1,
		"instead":        2,
		"intended":       1,
		"interface":      2,
		"internal":       1,
		"internally":     2,
		"is":             14,
		"it":             1,
		"just":           1,
		"last":           1,
		"less":           1,
		"license":        2,
		"like":           1,
		"loading":        1,
		"local":          1,
		"manu":           1,
		"martinez":       1,
		"matches":        2,
		"merged":         1,
		"method":         3,
		"methods":        2,
		"middleware":     6,
		"middlewares":    1,
		"mit":            1,
		"n":              2,
		"new":            2,
		"nolisting":      1,
		"non":            1,
		"not":            2,
		"obj":            1,
		"of":             6,
		"ones":           1,
		"operating":      1,
		"options":        5,
		"or":             3,
		"order":          1,
		"other":          1,
		"patch":          6,
		"path":           33,
		"pattern":        1,
		"permission":     1,
		"post":           6,
		"prefix":         2,
		"proxy":          1,
		"put":            6,
		"real":           1,
		"register":       1,
		"registers":      3,
		"relative":       16,
		"request":        1,
		"requests":       1,
		"reserved":       1,
		"reset":          1,
		"resources":      1,
		"respective":     1,
		"rest":           2,
		"return":         1,
		"returns":        1,
		"rights":         1,
		"root":           3,
		"route":          2,
		"router":         22,
		"routes":         5,
		"s":              2,
		"same":           1,
		"see":            2,
		"serve":          1,
		"server":         2,
		"serves":         1,
		"shared":         1,
		"shortcut":       8,
		"should":         4,
		"single":         3,
		"size":           1,
		"source":         1,
		"standardized":   1,
		"static":         14,
		"style":          1,
		"system":         4,
		"that":           5,
		"the":            18,
		"therefore":      1,
		"this":           2,
		"to":             6,
		"trace":          1,
		"url":            1,
		"usage":          1,
		"use":            7,
		"used":           6,
		"user":           1,
		"v":              4,
		"var":            1,
		"we":             1,
		"with":           3,
		"works":          1,
		"www":            1,
		"you":            1,
	}

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "", []byte(src), parser.ParseComments)

	samurai := NewSamuraiExtractor()
	ast.Walk(samurai, node)

	assert.NotEmpty(t, samurai.words)
	assert.Equal(t, len(expectedWords), len(samurai.words))
	for key, value := range expectedWords {
		assert.Equal(t, value, samurai.words[key], fmt.Sprintf("invalid number of occurrencies for element: %s", key))
	}
}

func TestFreqTable_OnNewlyCreatedSamuraiExtractor_ShouldReturnEmptyFreqTable(t *testing.T) {
	samurai := NewSamuraiExtractor()

	got := samurai.FreqTable()
	assert.Empty(t, got, fmt.Sprintf("frequency table should be empty: %v", got))
}

func TestFreqTable_OnSamuraiExtractorAfterExtraction_ShouldReturnFreqTableWithValues(t *testing.T) {
	src := `
		package main

		func main() {}
	`

	fs := token.NewFileSet()
	node, _ := parser.ParseFile(fs, "", []byte(src), parser.ParseComments)

	samurai := NewSamuraiExtractor()
	ast.Walk(samurai, node)

	freqTable := samurai.FreqTable()
	assert.NotEmpty(t, freqTable)
	assert.Equal(t, 1, len(freqTable))

	assert.Equal(t, 1, freqTable["main"], fmt.Sprintf("invalid number of occurrencies for element: main"))
}
