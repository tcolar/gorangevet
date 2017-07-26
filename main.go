package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"log"
	"os"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

type visitor struct {
	program *loader.Program
	failed  bool
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.RangeStmt: // check range statements
		key, _ := n.Key.(*ast.Ident)
		val, _ := n.Value.(*ast.Ident)
		if key == nil && val == nil { // nothing to check
			return nil
		}
		ast.Inspect(n.Body, func(n ast.Node) bool {
			u, ok := n.(*ast.UnaryExpr) // taking a pointer
			if !ok {
				return true
			}
			i, ok := u.X.(*ast.Ident) // to something
			if !ok {
				return true
			}
			if i.Name == key.Name { // that is the range key -> don't think that's a good idea
				fmt.Printf("%v Taking pointer to range key '%s'!\n",
					v.program.Fset.Position(i.Pos()), key)
				v.failed = true
			} else if i.Name == val.Name { // that is the range value -> don't think that's a good idea
				fmt.Printf("%v Taking pointer to range value '%s'!\n",
					v.program.Fset.Position(i.Pos()), val)
				v.failed = true
			}
			return true
		})
		return nil
	}

	return v
}

func main() {
	flag.Parse()
	imports := gotool.ImportPaths(flag.Args())
	if len(imports) == 0 {
		imports = []string{"."}
	}

	context := build.Default
	config := loader.Config{
		Build: &context,
	}
	_, err := config.FromArgs(imports, true)
	if err != nil {
		log.Fatalf("Failed to parse arguments: %s", err)
	}

	program, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load program code: %s", err)
	}

	failed := false
	for _, pkgInfo := range program.InitialPackages() {
		if pkgInfo.Pkg.Path() == "unsafe" {
			continue
		}
		v := &visitor{program: program}
		for _, f := range pkgInfo.Files {
			ast.Walk(v, f)
			failed = failed || v.failed
		}
	}
	if failed {
		os.Exit(1)
	}
	os.Exit(0)
}
