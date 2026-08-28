package highlight

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sitter "github.com/tree-sitter/go-tree-sitter/bindings/go"
)

type HighlightEngine struct {
	parser    *sitter.Parser
	wasmDir   string
	languages map[string]*sitter.Language
}

func NewHighlightEngine() (*HighlightEngine, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	wasmDir := filepath.Join(home, ".config", "moose", "tree-sitter")

	return &HighlightEngine{
		parser:    sitter.NewParser(),
		wasmDir:   wasmDir,
		languages: make(map[string]*sitter.Language),
	}, nil
}

func (h *HighlightEngine) LoadLanguage(ctx context.Context, langName string) (*sitter.Language, error) {
	if lang, exists := h.languages[langName]; exists {
		return lang, nil
	}

	wasmFile := fmt.Sprintf("tree-sitter-%s.wasm", langName)
	wasmPath := filepath.Join(h.wasmDir, wasmFile)

	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read grammar %s: %w", wasmPath, err)
	}

	lang, err := sitter.NewWasmLanguage(ctx, wasmBytes, fmt.Sprintf("tree_sitter_%s", langName))
	if err != nil {
		return nil, fmt.Errorf("failed to load WASM language %s: %w", langName, err)
	}

	h.languages[langName] = lang
	return lang, nil
}

func (h *HighlightEngine) ParseBuffer(ctx context.Context, langName string, content []byte) (*sitter.Tree, error) {
	lang, err := h.LoadLanguage(ctx, langName)
	if err != nil {
		return nil, err
	}

	h.parser.SetLanguage(lang)
	tree := h.parser.Parse(content, nil)
	return tree, nil
}
