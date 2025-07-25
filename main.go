package main

import (
	"flag"
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_enum.go")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of enumgen:\n")
	fmt.Fprintf(os.Stderr, "\tenumgen [flags] -type T [directory]\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	types := strings.Split(*typeNames, ",")

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var dir string
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}

	pkgs := loadPackages(args)
	sort.Slice(pkgs, func(i, j int) bool {
		iTest := strings.HasSuffix(pkgs[i].name, "_test")
		jTest := strings.HasSuffix(pkgs[j].name, "_test")
		if iTest && jTest {
			return !iTest
		}

		return len(pkgs[i].files) < len(pkgs[j].files)
	})

	generateAll(pkgs, types, dir)
}

func generateAll(pkgs []*Package, types []string, dir string) {
	for _, pkg := range pkgs {
		g := Generator{
			pkg: pkg,
		}

		g.Printf("// Code generated by \"enumgen %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
		g.Printf("\n")
		g.Printf("package %s", g.pkg.name)
		g.Printf("\n")

		var foundTypes, remainingTypes []string
		for _, typeName := range types {
			values := findValues(typeName, pkg)
			if len(values) > 0 {
				g.generate(typeName, values)
				foundTypes = append(foundTypes, typeName)
			} else {
				remainingTypes = append(remainingTypes, typeName)
			}
		}

		if len(foundTypes) == 0 {
			continue
		}

		if len(remainingTypes) > 0 && output != nil && *output != "" {
			log.Fatalf("cannot write to single file (-output=%q) when matching types are found in multiple packages", *output)
		}
		types = remainingTypes

		src := g.format()

		outputName := *output
		if outputName == "" {
			outputName = filepath.Join(dir, baseName(pkg, foundTypes[0]))
		}

		err := os.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}

	if len(types) > 0 {
		log.Fatalf("no values defined for types: %s", strings.Join(types, " "))
	}
}

func baseName(pkg *Package, typename string) string {
	suffix := "enum.go"
	if pkg.hasTestFiles {
		suffix = "enum_test.go"
	}

	return fmt.Sprintf("%s_%s", strings.ToLower(typename), suffix)
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func loadPackages(patterns []string) []*Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedModule,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}

	if len(pkgs) == 0 {
		log.Fatalf("error: no packages matching %v", strings.Join(patterns, " "))
	}

	out := make([]*Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		p := &Package{
			name:  pkg.Name,
			defs:  pkg.TypesInfo.Defs,
			files: make([]*File, 0, len(pkg.Syntax)),
		}

		for _, file := range pkg.Syntax {
			p.files = append(p.files, &File{
				file: file,
				pkg:  p,
			})
		}

		for _, f := range pkg.GoFiles {
			if strings.HasSuffix(f, "_test.go") {
				p.hasTestFiles = true
				break
			}
		}

		out = append(out, p)
	}

	return out
}

func findValues(typeName string, pkg *Package) []Value {
	values := make([]Value, 0, 100)
	for _, file := range pkg.files {
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	return values
}
