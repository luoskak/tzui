package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	root, _ := os.Getwd()
	filePath := filepath.Join(
		root,
		"..",
		"..",
		"model",
		"user_log",
		"interface.go",
	)
	f, err := parser.ParseFile(token.NewFileSet(), filePath, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	for _, fdcs := range f.Decls {
		switch decl := fdcs.(type) {
		case *ast.GenDecl:
			for _, s := range decl.Specs {
				switch spec := s.(type) {
				case *ast.TypeSpec:
					switch e := spec.Type.(type) {
					case *ast.StructType:
						fmt.Println(spec)
						fmt.Println(e)
						fmt.Println(e)

					}
				}
			}
		}
	}
}
