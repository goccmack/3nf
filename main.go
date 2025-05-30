package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccmack/3nf/ast"
	"github.com/goccmack/3nf/gen/dot"
	"github.com/goccmack/3nf/gen/sql"
	"github.com/goccmack/3nf/lexer"
	"github.com/goccmack/3nf/parser"
)

var (
	help    = flag.Bool("h", false, "Help")
	ermFile = ""
	ermDir  = ""
)

func main() {
	getParams()
	lex := lexer.NewFile(ermFile)
	pf, errs := parser.Parse(lex)
	if len(errs) != 0 {
		parseErrors(ermFile, errs)
	}
	schema := ast.Build(pf, lex, ermFile)
	schema.Check()
	sql.Gen(ermDir, schema)
	dot.Gen(ermDir, schema)
}

func getParams() {
	flag.Parse()
	if *help {
		fmt.Println(usage)
		os.Exit(0)
	}
	if flag.NArg() != 1 {
		fmt.Println("Error: exactly 1 data model must be specified")
		fmt.Print(usage)
		os.Exit(1)
	}
	ermFile = flag.Arg(0)
	ermDir, _ = filepath.Split(ermFile)
}

func parseErrors(fname string, errs []*parser.Error) {
	fmt.Println("Parse Errors:")
	ln := errs[0].Line
	for _, err := range errs {
		if err.Line == ln {
			fmt.Printf("%s:%d:%d: %s",
				fname, err.Line, err.Column, err)
		}
	}
	os.Exit(1)
}

const usage = `use 3nf <path to model file>`
