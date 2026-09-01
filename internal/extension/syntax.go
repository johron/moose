package extension

import (
	lua "github.com/yuin/gopher-lua"
	"moose/internal/extension/highlight"
)

func GetSyntaxTable(em *ExtensionManager) *lua.LTable {
	syntax := em.L.NewTable()
	em.L.SetFuncs(syntax, map[string]lua.LGFunction{
		"registerLang": luaRegisterLanguage,
	})
	return syntax
}

func luaRegisterLanguage(L *lua.LState) int {
	idx := 1
	if L.GetTop() >= 2 && L.Get(1).Type() == lua.LTTable {
		tb1 := L.ToTable(1)
		if tb1.RawGetString("name") == lua.LNil {
			idx = 2
		}
	}

	tb := L.CheckTable(idx)

	langName := tb.RawGetString("name").String()
	parserPath := tb.RawGetString("parser_path").String()
	queryScm := tb.RawGetString("query").String()

	if langName == "" || langName == "nil" || parserPath == "" || parserPath == "nil" {
		L.ArgError(idx, "expected table with 'name' and 'parser_path' fields")
		return 0
	}

	err := highlight.RegisterTreeSitterLang(langName, parserPath, queryScm)
	if err != nil {
		L.RaiseError("failed to register language %s: %v", langName, err)
		return 0
	}

	return 0
}