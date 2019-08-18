package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/eroatta/token/amap"
	"github.com/eroatta/token/greedy"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/samurai"
	"gopkg.in/src-d/go-git.v4"

	"github.com/eroatta/token/conserv"

	"github.com/eroatta/src-reader/extractors"

	"github.com/eroatta/src-reader/repositories"
)

func main() {
	newGoodMain("https://github.com/src-d/go-siva")
}

func newGoodMain(url string) {
	// stage: clone (and retrieve files)
	files, err := clone(url)
	if err != nil {
		log.Fatalf("Error reading repository %s: %v", url, err)
	}

	// parse (generate ASTs)
	trees := make([]*ast.File, 0)
	astc := parse(files)
	for ast := range astc {
		trees = append(trees, ast)
	}

	// TODO: merge ASTs considering packages

	// for each package
	// apply the set of miners (preprocessors)
	occurc, _ := mine(trees)
	freqTable := samurai.NewFrequencyTable()
	for input := range occurc {
		for token, count := range input {
			//TODO: review
			if len(token) == 1 {
				continue
			}
			freqTable.SetOccurrences(token, count)
		}
	}
	fmt.Println(freqTable)

	// for each package (AST)
	// apply the set of splitters + expanders
	//numberProcessors := 1
	identc := make(chan identifier)
	go func() {
		for _, t := range trees {
			ast.Inspect(t, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					ident := identifier{
						file:       "main.go",
						position:   fn.Pos(),
						name:       fn.Name.String(),
						typ:        "FuncDecl",
						splits:     make(map[string][]string),
						expansions: make(map[string][]string),
					}
					identc <- ident
				}

				return true
			})
		}
		close(identc)
	}()

	log.Println("Closed identc...")

	tCtx := samurai.NewTokenContext(freqTable, freqTable)
	splittedc := split(identc, tCtx)
	expandedc := expand(splittedc)

	for ident := range expandedc {
		log.Println("Identifier received")
		for alg, splits := range ident.splits {
			log.Println(fmt.Sprintf("FuncDecl \"%s\" Splitted into: %v by %s", ident.name, splits, alg))
		}

		for alg, expans := range ident.expansions {
			log.Println(fmt.Sprintf("FuncDecl \"%s\" Expanded into: %v by %s", ident.name, expans, alg))
		}
	}
}

type identifier struct {
	file       string
	position   token.Pos
	name       string
	typ        string
	node       *ast.Node
	splits     map[string][]string
	expansions map[string][]string
}

func clone(url string) (<-chan file, error) {
	/*
		CLONING STEP: retrieves the source code from Github.
		It validates the given Github URI, clones the repository, checks for errors
		and also extract the list of file names that belong to the project.
		Input: Github URI.
		Output: a list of files and a filesystem so we can read them.
	*/

	// TODO: rename package to: "repository"
	repository, err := repositories.Clone(repositories.GoGitClonerFunc, url)
	if err != nil {
		return nil, err
	}

	files, err := repositories.Filenames(repository)
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		for _, f := range files {
			if !strings.HasSuffix(f, ".go") {
				continue
			}
			out <- f
		}
		close(out)
	}()

	return retrieve(repository, out), nil
}

type file struct {
	name string
	raw  []byte
}

func retrieve(repo *git.Repository, namesc <-chan string) chan file {
	filesc := make(chan file)
	go func() {
		for n := range namesc {
			rawFile, err := repositories.File(repo, n)
			// TODO: review errors (do I need error channel?)
			if err != nil {
				continue
			}

			file := file{
				name: n,
				raw:  rawFile,
			}
			filesc <- file
		}
		close(filesc)
	}()

	return filesc
}

func parse(filesc <-chan file) chan *ast.File {
	/*
		PARSING STEP: parses a file and creates and AST so it can be explored lately.
		It filters non-go files and skips them, so we can focus only on Go source code.
		It also retrieves the binary files and parses it, using the ast/parser from Go.
		Input: a list of files and a filesystem so we can read them.
		Output: a list of ASTs, where each AST represents a file.
		Improvement opportunity: we could merge the ASTs into one collection of ast.Package.
	*/
	fset := token.NewFileSet()

	astc := make(chan *ast.File)
	go func() {
		for file := range filesc {
			node, err := parser.ParseFile(fset, file.name, file.raw, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
				continue
			}
			astc <- node
		}
		close(astc)
	}()

	return astc
}

func mine(trees []*ast.File) (chan map[string]int, chan map[string]amap.TokenScope) {
	/*
		MINING STEP: extracts all the information required by each splitting and expanding algorithm.
		It traverses each AST with every defined miner, so they can extract the required info.
		Input: a list of *ast.File nodes, and a list of Miners (new interface?)
		Output: ?
	*/
	occurc := make(chan map[string]int)
	go func() {
		for _, t := range trees {
			ext := extractors.NewSamuraiExtractor()
			ast.Walk(ext, t)
			occurc <- ext.FreqTable()
		}
		close(occurc)
	}()

	scopesc := make(chan map[string]amap.TokenScope)
	go func() {
		for _, t := range trees {
			ext := extractors.NewAmap("main.go")
			ast.Walk(ext, t)
			scopesc <- ext.Scopes()
		}
		close(scopesc)
	}()

	return occurc, scopesc
}

func split(identc <-chan identifier, tCtx samurai.TokenContext) chan identifier {
	/*
		SPLITTING STEP: it splits all the identifiers in a AST, applying the given set of Splitters.
		Input: a list of *ast.File nodes, and a list of Splitters.
		Output: ?
	*/
	splittedc := make(chan identifier)
	go func() {
		for ident := range identc {
			log.Println("Splitting...")
			ident.splits["conserv"] = conserv.Split(ident.name)
			ident.splits["greedy"] = greedy.Split(ident.name, greedy.DefaultList)
			ident.splits["samurai"] = samurai.Split(ident.name, tCtx, lists.Prefixes, lists.Suffixes)

			splittedc <- ident
		}
		close(splittedc)
	}()

	return splittedc
}

func expand(identc <-chan identifier) chan identifier {
	expandedc := make(chan identifier)
	go func() {
		for ident := range identc {
			log.Println("Expanding...")
			// TODO: add expansions

			expandedc <- ident
		}

		close(expandedc)
	}()

	return expandedc
}
