package highlight

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>

// Helper to resolve the C function pointer for tree_sitter_<lang>()
void* load_symbol(void* handle, const char* symbol) {
    return dlsym(handle, symbol);
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
	"strings"
	"unicode"
	"path/filepath"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

type LanguageSpec struct {
	Name     string
	Language *sitter.Language
	Query    *sitter.Query
}

type Registry struct {
	mu        sync.RWMutex
	languages map[string]*LanguageSpec
}

var GlobalRegistry = &Registry{
	languages: make(map[string]*LanguageSpec),
}

func CleanScmQuery(scm string) string {
	cleaned := strings.Map(func(r rune) rune {
		if r == '\uFEFF' { // Remove UTF-8 Byte Order Mark (BOM)
			return -1
		}
		if unicode.IsSpace(r) && r != '\n' && r != '\r' && r != '\t' {
			return ' '
		}
		return r
	}, scm)

	return strings.TrimSpace(cleaned)
}

func RegisterTreeSitterLang(langName, parserPath, queryScm string) error {
	absPath, err := filepath.Abs(parserPath)
	if err == nil {
		parserPath = absPath
	}

	cPath := C.CString(parserPath)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.dlopen(cPath, C.RTLD_NOW)
	if handle == nil {
		errStr := C.GoString(C.dlerror())
		return fmt.Errorf("failed to load parser library at %s: %s", parserPath, errStr)
	}

	symbolName := C.CString(fmt.Sprintf("tree_sitter_%s", langName))
	defer C.free(unsafe.Pointer(symbolName))

	sym := C.load_symbol(handle, symbolName)
	if sym == nil {
		return fmt.Errorf("symbol 'tree_sitter_%s' not found in %s", langName, parserPath)
	}

	ptr := unsafe.Pointer(sym)
	sitterLang := sitter.NewLanguage(ptr)
	if sitterLang == nil {
		return fmt.Errorf("invalid TSLanguage pointer returned for %s", langName)
	}

	cleanScm := CleanScmQuery(queryScm)
	query, err := sitter.NewQuery(sitterLang, string([]byte(cleanScm)))
	if err != nil {
		return fmt.Errorf("failed to compile highlight query for %s: %w", langName, err)
	}

	GlobalRegistry.mu.Lock()
	GlobalRegistry.languages[langName] = &LanguageSpec{
		Name:     langName,
		Language: sitterLang,
		Query:    query,
	}
	GlobalRegistry.mu.Unlock()

	return nil
}

func (r *Registry) GetLanguage(langName string) (*LanguageSpec, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	spec, ok := r.languages[langName]
	return spec, ok
}