package extension

import (
	"embed"
	"fmt"
	"moose/internal/editor"
	"moose/internal/extension/highlight"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type ExtensionManager struct {
	L           *lua.LState
	M           *editor.Model
	HE			*highlight.HighlightEngine
	LoadedFiles []string

	currentDiskDir  string
	currentEmbedDir string
}

//go:embed lua
var embeddedScripts embed.FS

func NewExtensionManager(m *editor.Model) *ExtensionManager {
	em := &ExtensionManager{
		L: lua.NewState(),
		M: m,
		HE: highlight.NewHighlightEngine(),
	}

	lua.OpenPackage(em.L)

	em.registerAPI()
	em.registerExtensionSearcher()

	if err := em.LoadEmbeddedFile("init.lua"); err != nil {
		m.Mode = editor.ModeNormal
		m.BM.PaletteBuffer.Clear()
		m.BM.PaletteBuffer.Insert("moose.error:Lua error " + err.Error())
	}

	if err := em.LoadFile("/home/johron/.config/moose/moose.lua"); err != nil {
		m.Mode = editor.ModeNormal
		m.BM.PaletteBuffer.Clear()
		m.BM.PaletteBuffer.Insert("moose.error:Lua error " + err.Error())
	}

	return em
}

func (em *ExtensionManager) registerExtensionSearcher() {
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

	extensionSearcher := em.L.NewFunction(func(L *lua.LState) int {
		modName := L.CheckString(1)
		fileName := strings.ReplaceAll(modName, ".", "/") + ".lua"
		var errorsLogged []string

		if em.currentDiskDir != "" {
			targetDiskFile := filepath.Join(em.currentDiskDir, fileName)
			if fn, loadErr := L.LoadFile(targetDiskFile); loadErr == nil {
				return pushAndReturn(L, fn)
			} else {
				errorsLogged = append(errorsLogged, fmt.Sprintf("no relative file: %s", targetDiskFile))
			}
		}

		if em.currentEmbedDir != "" {
			targetEmbedFile := filepath.Join(em.currentEmbedDir, fileName)
			if bytes, err := embeddedScripts.ReadFile(targetEmbedFile); err == nil {
				if fn, err := L.LoadString(string(bytes)); err == nil {
					return pushAndReturn(L, fn)
				} else {
					L.RaiseError("failed to compile embedded module %s: %v", modName, err)
					return 0
				}
			} else {
				errorsLogged = append(errorsLogged, fmt.Sprintf("no relative embed file: %s", targetEmbedFile))
			}
		}

		L.Push(lua.LString("\n\t" + strings.Join(errorsLogged, "\n\t")))
		return 1
	})

	idx2 := packageLoaders.RawGetInt(2)
	idx3 := packageLoaders.RawGetInt(3)
	idx4 := packageLoaders.RawGetInt(4)

	packageLoaders.RawSetInt(2, extensionSearcher)
	packageLoaders.RawSetInt(3, idx2)
	packageLoaders.RawSetInt(4, idx3)
	if idx4 != lua.LNil {
		packageLoaders.RawSetInt(5, idx4)
	}
}

func pushAndReturn(L *lua.LState, fn *lua.LFunction) int {
	L.Push(fn)
	return 1
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
	moose.RawSetString("config", GetConfigTable(em))
}

func (em *ExtensionManager) LoadFile(path string) error {
	oldDiskDir := em.currentDiskDir
	em.currentDiskDir = filepath.Dir(path)
	defer func() { em.currentDiskDir = oldDiskDir }()

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

	oldEmbedDir := em.currentEmbedDir
	em.currentEmbedDir = filepath.Dir(embedPath)
	defer func() { em.currentEmbedDir = oldEmbedDir }()

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
