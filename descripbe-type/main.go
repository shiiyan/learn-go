package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

func describeType(t types.Type) {
	switch x := t.(type) {
	case *types.Named:
		fmt.Printf("Named type: %s\n", x.Obj().Name())
		fmt.Printf("Underlying: %s\n", x.Underlying().String())
	case *types.Basic:
		// A basic built-in type like “int”
		fmt.Printf("Basic type: %s\n", x.Name())
	default:
		fmt.Printf("Other type: %T -> %s\n", x, x.String())
	}
}

func main() {
	// 1. Create a FileSet and parse a Go source file (in this example, we embed the code as a string).
	src := `
        package main

        type MyInt int
        type MyInt2 = MyInt
    `

	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, "example.go", src, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	// 2. Prepare type‐checking
	conf := types.Config{
		// If you need to resolve imports, set Up an Importer (e.g. importer.Default()).
		// For a single file with no imports, the zero‐value importer is fine.
	}
	info := &types.Info{
		// We only need TypeOf and Defs in this simple example:
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
	}

	// 3. Type‐check the “example.go” file as-package “main”
	pkg, err := conf.Check("main", fset, []*ast.File{fileNode}, info)
	if err != nil {
		panic(err)
	}

	// 4. Look up “MyInt” and “MyInt2” in the package’s scope:
	scope := pkg.Scope()

	objMyInt := scope.Lookup("MyInt")
	if objMyInt == nil {
		panic("could not find MyInt")
	}
	tyMyInt := objMyInt.Type() // this is a types.Type
	fmt.Println("---- MyInt: ----")
	describeType(tyMyInt)

	objMyInt2 := scope.Lookup("MyInt2")
	if objMyInt2 == nil {
		panic("could not find MyInt2")
	}
	tyMyInt2 := objMyInt2.Type()
	fmt.Println("---- MyInt2: ----")
	describeType(tyMyInt2)
}
