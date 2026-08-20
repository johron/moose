package extension

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"moose/internal/editor"
)

type ExtensionManager struct {
	L *lua.LState
	Model *editor.Model
	LoadedFiles []string
}

func NewExtensionManager(m *editor.Model) *ExtensionManager {
	em := &ExtensionManager{
		L: lua.NewState(),
		Model: m,
	}

	em.registerAPI()
	return em
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