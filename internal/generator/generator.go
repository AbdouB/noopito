package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// GenerateNoOp generates a no-op implementation of interfaceName from pkgPath
// and writes it to outputDir.
func GenerateNoOp(pkgPath, interfaceName, outputDir string) error {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedName | packages.NeedModule}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return fmt.Errorf("loading package: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return fmt.Errorf("package has errors")
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("package %s not found", pkgPath)
	}
	pkg := pkgs[0]
	obj := pkg.Types.Scope().Lookup(interfaceName)
	if obj == nil {
		return fmt.Errorf("interface %s not found in package %s", interfaceName, pkgPath)
	}
	named, ok := obj.(*types.TypeName)
	if !ok {
		return fmt.Errorf("%s is not a type", interfaceName)
	}
	iface, ok := named.Type().Underlying().(*types.Interface)
	if !ok {
		return fmt.Errorf("%s is not an interface", interfaceName)
	}
	iface = iface.Complete()

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "package %s\n\n", pkg.Name)
	implName := interfaceName + "NoOp"
	fmt.Fprintf(&buf, "type %s struct{}\n\n", implName)

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		sig := m.Type().(*types.Signature)

		params := make([]string, 0, sig.Params().Len())
		for j := 0; j < sig.Params().Len(); j++ {
			p := sig.Params().At(j)
			name := p.Name()
			if name == "" {
				name = fmt.Sprintf("arg%d", j)
			}
			params = append(params, fmt.Sprintf("%s %s", name, typeString(p.Type())))
		}
		paramsStr := strings.Join(params, ", ")

		results := make([]string, 0, sig.Results().Len())
		for j := 0; j < sig.Results().Len(); j++ {
			r := sig.Results().At(j)
			tstr := typeString(r.Type())
			if r.Name() != "" {
				results = append(results, fmt.Sprintf("%s %s", r.Name(), tstr))
			} else {
				results = append(results, tstr)
			}
		}
		resultsStr := strings.Join(results, ", ")
		if len(results) > 1 {
			resultsStr = "(" + resultsStr + ")"
		}

		fmt.Fprintf(&buf, "func (n %s) %s(%s) %s {\n", implName, m.Name(), paramsStr, resultsStr)
		if sig.Results().Len() > 0 {
			for j := 0; j < sig.Results().Len(); j++ {
				r := sig.Results().At(j)
				fmt.Fprintf(&buf, "\tvar r%d %s\n", j, typeString(r.Type()))
			}
			buf.WriteString("\treturn ")
			for j := 0; j < sig.Results().Len(); j++ {
				if j > 0 {
					buf.WriteString(", ")
				}
				fmt.Fprintf(&buf, "r%d", j)
			}
			buf.WriteString("\n")
		}
		buf.WriteString("}\n\n")
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting code: %w", err)
	}

	outPath := filepath.Join(outputDir, fmt.Sprintf("%s_noop.go", strings.ToLower(interfaceName)))
	return os.WriteFile(outPath, formatted, 0o644)
}

func typeString(t types.Type) string {
	return types.TypeString(t, func(p *types.Package) string {
		if p == nil {
			return ""
		}
		return p.Name()
	})
}
