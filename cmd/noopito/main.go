package main

import (
	"flag"
	"log"

	"github.com/AbdouB/noopito/internal/generator"
)

func main() {
	pkg := flag.String("package", "", "package import path")
	iface := flag.String("interface", "", "interface name")
	outDir := flag.String("output", ".", "output directory")
	flag.Parse()

	if *pkg == "" || *iface == "" {
		log.Fatal("package and interface flags are required")
	}

	if err := generator.GenerateNoOp(*pkg, *iface, *outDir); err != nil {
		log.Fatal(err)
	}
}
