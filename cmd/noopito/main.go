package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/AbdouB/noopito/internal/generator"
)

func main() {
	pkg := flag.String("package", "", "package import path (defaults to current package)")
	iface := flag.String("interface", "", "interface name")
	outDir := flag.String("output", ".", "output directory")
	flag.Parse()

	if *iface == "" {
		log.Fatal("interface flag is required")
	}

	pkgPath := *pkg
	if pkgPath == "" {
		var err error
		pkgPath, err = currentPackage()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := generator.GenerateNoOp(pkgPath, *iface, *outDir); err != nil {
		log.Fatal(err)
	}
}

func currentPackage() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	gomod := os.Getenv("GOMOD")
	if gomod == "" {
		return ".", nil
	}

	data, err := os.ReadFile(gomod)
	if err != nil {
		return "", err
	}

	mf, err := modfile.Parse(gomod, data, nil)
	if err != nil {
		return "", err
	}

	modDir := filepath.Dir(gomod)
	rel, err := filepath.Rel(modDir, wd)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	if rel == "." {
		return mf.Module.Mod.Path, nil
	}
	return path.Join(mf.Module.Mod.Path, rel), nil
}
