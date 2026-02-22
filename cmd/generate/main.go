// Copyright 2026 Bohdan Shtepan.
// Licensed under the MIT License.

// Command generate reads the LSP metaModel.json specification and produces
// Go source files for the protocol package.
//
// Usage:
//
//	go run github.com/modern-dev/go-lsp/cmd/generate [-o dir] [-model path] [-ref tag]
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/modern-dev/go-lsp/internal/generate"
)

const defaultRef = "release/protocol/3.17.6-next.14"

var httpClient = &http.Client{ //nolint:exhaustruct,gochecknoglobals
	Timeout: 30 * time.Second, //nolint:mnd
}

func main() {
	outDir := flag.String("o", "protocol", "Output directory for generated files")
	modelPath := flag.String("model", "", "Path to a local metaModel.json (skips download)")
	ref := flag.String("ref", defaultRef, "Git ref / tag to fetch metaModel.json from")

	flag.Parse()

	data, err := loadModel(*modelPath, *ref)
	if err != nil {
		log.Fatalf("load model: %v", err)
	}

	var model generate.Model
	if err := json.Unmarshal(data, &model); err != nil { //nolint:noinlineerr
		log.Fatalf("parse metaModel.json: %v", err)
	}

	fmt.Printf("LSP version: %s\n", model.MetaData.Version)
	fmt.Printf("Structures:    %d\n", len(model.Structures))
	fmt.Printf("Enumerations:  %d\n", len(model.Enumerations))
	fmt.Printf("TypeAliases:   %d\n", len(model.TypeAliases))
	fmt.Printf("Requests:      %d\n", len(model.Requests))
	fmt.Printf("Notifications: %d\n", len(model.Notifications))

	gen := generate.NewGenerator(&model)

	out, err := gen.Generate()
	if err != nil {
		log.Fatalf("generate: %v", err)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil { //nolint:gosec,mnd,noinlineerr
		log.Fatalf("mkdir %s: %v", *outDir, err)
	}

	type namedFile struct {
		name    string
		content []byte
	}

	files := []namedFile{
		{"types_gen.go", out.Types},
		{"server_gen.go", out.Server},
		{"client_gen.go", out.Client},
	}

	for _, fil := range files {
		path := filepath.Join(*outDir, fil.name)
		if err := os.WriteFile( //nolint:gosec,noinlineerr
			path,
			fil.content,
			0o644, //nolint:mnd
		); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}

		fmt.Printf("Wrote %s (%d bytes)\n", path, len(fil.content))
	}
}

// loadModel returns the raw bytes of metaModel.json, either from a local file
// or by downloading it from the vscode-languageserver-node repository.
func loadModel(localPath, ref string) ([]byte, error) {
	if localPath != "" {
		fmt.Printf("Reading local model: %s\n", localPath)

		return os.ReadFile(filepath.Clean(localPath)) //nolint:wrapcheck
	}

	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/%s/protocol/metaModel.json",
		ref,
	)

	fmt.Printf("Downloading metaModel.json from %s\n", url)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url) //nolint:err113
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return data, nil
}
