package highlight

import (
	"fmt"

	tcell "github.com/gdamore/tcell/v3"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tspack "github.com/xberg-io/tree-sitter-language-pack/packages/go"
)

type StyleMap map[string]tcell.Style

type Highlighter struct {
	parser *tree_sitter.Parser
	theme  StyleMap
}

func NewHighlighter(theme StyleMap) (*Highlighter, error) {
	parser := tree_sitter.NewParser()

	return &Highlighter{
		parser: parser,
		theme:  theme,
	}, nil
}

func (h *Highlighter) Close() {
	if h.parser != nil {
		h.parser.Close()
	}
}

func (h *Highlighter) HighlightBuffer(langName string, code []byte, highlightScm []byte) ([]HighlightSpan, *tree_sitter.Tree, error) {
	lang, err := tspack.GetLanguage(langName)
	if err != nil || lang == nil {
		return nil, nil, fmt.Errorf("could not load language %s: %w", langName, err)
	}

	if err := h.parser.SetLanguage(lang); err != nil {
		return nil, nil, fmt.Errorf("failed to set language: %w", err)
	}

	tree := h.parser.Parse(code, nil)
	if tree == nil {
		return nil, nil, fmt.Errorf("failed to parse document")
	}

	query, err := tree_sitter.NewQuery(lang, string(highlightScm))
	if err != nil {
		tree.Close()
		return nil, nil, fmt.Errorf("invalid query file: %w", err)
	}
	defer query.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	captures := cursor.Captures(query, tree.RootNode(), code)
	var spans []HighlightSpan

	for match, captureIdx := captures.Next(); match != nil; match, captureIdx = captures.Next() {
		capture := match.Captures[captureIdx]
		node := capture.Node
		captureName := query.CaptureNames()[capture.Index]

		if style, ok := h.theme[captureName]; ok {
			spans = append(spans, HighlightSpan{
				StartByte: uint32(node.StartByte()),
				EndByte:   uint32(node.EndByte()),
				Style:     style,
			})
		}
	}

	return spans, tree, nil
}

type HighlightSpan struct {
	StartByte uint32
	EndByte   uint32
	Style     tcell.Style
}