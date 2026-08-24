package extension

import (
	lua "github.com/yuin/gopher-lua"
)

func GetConfigTable(em *ExtensionManager) *lua.LTable {
	config := em.L.NewTable()
	em.L.SetFuncs(config, map[string]lua.LGFunction{
		"set": func(L *lua.LState) int {
			em.M.DebugLog += "hi from lua"
			return 1
		},
	})

	return config
}
