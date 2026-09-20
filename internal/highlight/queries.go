package highlight

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
)

//go:embed queries/*.scm
var QueryFS embed.FS

func GetQuery(lang string) ([]byte, error) {
	path := fmt.Sprintf("queries/%s.scm", lang)
	data, err := QueryFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no highlight query found for language %q: %w", lang, err)
	}
	return data, nil
}

func DetectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext { // TODO: set this in configuration
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".js", ".jsx":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".c", ".h":
		return "c"
	case ".cpp", ".hpp", ".cc":
		return "cpp"
	case ".json":
		return "json"
	case  ".lua":
		return "lua"
	default:
		return "" // Unknown/Plaintext
	}
}