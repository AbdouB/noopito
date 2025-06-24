# noopito

`noopito` is a small CLI utility that generates no-op implementations for Go interfaces. The tool is designed to work with `go generate` and is useful when you need stub implementations for testing or scaffolding.

## Installation

```bash
go install github.com/AbdouB/noopito/cmd/noopito@latest
```

Run the tool without installing it using `go run`:

```bash
go run github.com/AbdouB/noopito/cmd/noopito@latest \
  -package=your/module/pkg -interface=InterfaceName -output=./mocks
```

## Usage

Run the tool by specifying the package path and the interface name. If the `-package` flag is omitted when running via `go generate`, the current package is used automatically. By default, the generated file is written to the current directory.

```bash
noopito -package=your/module/pkg -interface=InterfaceName -output=./mocks
```

This will create `interfacename_noop.go` in the output directory containing a struct named `InterfaceNameNoOp` that implements all methods with zero-value returns.

### Using `go generate`

Add a `//go:generate` directive near the interface definition:

```go
//go:generate noopito -interface=InterfaceName -output=./mocks
```

Then run `go generate ./...` to produce the file automatically.

## License

MIT
