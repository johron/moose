package extension

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"moose/internal/editor"
	"embed"
	"path/filepath"
	"strings"
)

type ExtensionManager struct {
	L *lua.LState
	Model *editor.Model
	LoadedFiles []string
}

//go:embed lua
var embeddedScripts embed.FS

func NewExtensionManager(m *editor.Model) *ExtensionManager {
	em := &ExtensionManager{
		L: lua.NewState(),
		Model: m,
	}

	em.registerAPI()
	em.registerEmbedSearcher()

	if err := em.LoadEmbeddedFile("init.lua"); err != nil {
		m.Mode = editor.ModeNormal
		m.BM.PaletteBuffer.Clear()
		m.BM.PaletteBuffer.Insert("moose.error:Lua error " + err.Error())
	}

	return em
}

func (em *ExtensionManager) registerEmbedSearcher() {
	if em.L == nil {
		return
	}

	pkg := em.L.GetGlobal("package")
	if pkg.Type() == lua.LTNil {
		return
	}

	loadersVal := em.L.GetField(pkg, "loaders")
	packageLoaders, ok := loadersVal.(*lua.LTable)
	if !ok {
		return
	}

	embedSearcher := em.L.NewFunction(func(L *lua.LState) int {
		modName := L.CheckString(1)

		fileName := strings.ReplaceAll(modName, ".", "/") + ".lua"
		embedPath := filepath.Join("lua", fileName)

		bytes, err := embeddedScripts.ReadFile(embedPath)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("\n\tno embedded file: %s", embedPath)))
			return 1
		}

		fn, err := L.LoadString(string(bytes))
		if err != nil {
			L.RaiseError("failed to compile embedded module %s: %v", modName, err)
			return 0
		}

		L.Push(fn)
		return 1
	})

	packageLoaders.Append(embedSearcher)
}


func (em *ExtensionManager) Close() {
	if em.L != nil {
		em.L.Close()
		em.L = nil
	}
}

func (em *ExtensionManager) registerAPI() {
	if em.L == nil {
		return
	}

	moose := em.L.NewTable()
	em.L.SetGlobal("ms", moose)
}

func (em *ExtensionManager) LoadFile(path string) error {
	if err := em.L.DoFile(path); err != nil {
		return err
	}

	em.LoadedFiles = append(em.LoadedFiles, path)
	return nil
}

func (em *ExtensionManager) LoadEmbeddedFile(path string) error {
	if em.L == nil {
		return fmt.Errorf("lua state is closed")
	}

	embedPath := filepath.Join("lua", path)

	bytes, err := embeddedScripts.ReadFile(embedPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded script %s: %w", path, err)
	}

	if err := em.L.DoString(string(bytes)); err != nil {
		return err
	}

	em.LoadedFiles = append(em.LoadedFiles, path)
	return nil
}

func (em *ExtensionManager) LoadString(name string, src string) error {
	if em.L == nil {
		return fmt.Errorf("lua state is closed")
	}

	if err := em.L.DoString(src); err != nil {
		return err
	}

	em.LoadedFiles = append(em.LoadedFiles, name)
	return nil
}