package extension

import (
	"fmt"
	"reflect"
	"strings"

	"moose/internal/editor"

	"github.com/creasty/defaults"
	"github.com/gdamore/tcell/v3"
	lua "github.com/yuin/gopher-lua"
)

type Config struct {
	StyleDefault tcell.Style
	Colors       Colors
	Properties   Properties
}

type Properties struct {
	GutterWidth    int  `default:"4"`
	TabSpaces      bool `default:"false"`
	TabWidthSpaces int  `default:"4"`
}

type Colors struct {
	MainBackground            string `default:"#0E142E"`
	MainForeground            string `default:"#dddddd"`
	LineNumberBackground      string `default:"#090603"`
	LineNumberForeground      string `default:"#dddddd"`
	CursorColor               string `default:"#777777"`
	CursorColorWrite          string `default:"#dddddd"`
	PaletteBarBackground      string `default:"#3b3b3b"`
	PaletteBarForeground      string `default:"#dddddd"`
	PaletteInputBackground    string `default:"#090603"`
	PaletteInputForeground    string `default:"#dddddd"`
	InfoMsgBackground         string `default:"#090603"`
	InfoMsgForeground         string `default:"#dddddd"`
	WarnMsgBackground         string `default:"#090603"`
	WarnMsgForeground         string `default:"#efc541"`
	ErrorMsgBackground        string `default:"#090603"`
	ErrorMsgForeground        string `default:"#ff5e56"`
	WorkspaceBackground       string `default:"#3b3b3b"`
	WorkspaceForeground       string `default:"#b9b9b9"`
	WorkspaceBackgroundActive string `default:"#dddddd"`
	WorkspaceForegroundActive string `default:"#3b3b3b"`
	TabBackground             string `default:"#3b3b3b"`
	TabForeground             string `default:"#b9b9b9"`
	TabBackgroundActive       string `default:"#3b3b3b"`
	TabForegroundActive       string `default:"#dddddd"`
}

func DefaultConfig() Config {
	config := Config{}
	if err := defaults.Set(&config); err != nil {
		panic(err)
	}

	config.StyleDefault = tcell.StyleDefault.Background(tcell.GetColor(config.Colors.MainBackground)).Foreground(tcell.GetColor(config.Colors.MainForeground))

	return config
}

func snakeToCamel(s string) string {
	parts := strings.Split(strings.ToLower(s), "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func parseStructTable(dst any, table *lua.LTable) error {
	if err := defaults.Set(dst); err != nil {
		return err
	}

	val := reflect.ValueOf(dst)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst must be a pointer to struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		luaValue := table.RawGetString(fieldType.Name)
		if luaValue == lua.LNil {
			continue
		}

		switch fieldVal.Kind() {
		case reflect.String:
			if v, ok := luaValue.(lua.LString); ok {
				fieldVal.SetString(string(v))
			}

		case reflect.Bool:
			if v, ok := luaValue.(lua.LBool); ok {
				fieldVal.SetBool(bool(v))
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, ok := luaValue.(lua.LNumber); ok {
				fieldVal.SetInt(int64(v))
			}

		case reflect.Float32, reflect.Float64:
			if v, ok := luaValue.(lua.LNumber); ok {
				fieldVal.SetFloat(float64(v))
			}

		case reflect.Struct:
			if tbl, ok := luaValue.(*lua.LTable); ok {
				nested := fieldVal.Addr().Interface()
				if err := parseStructTable(nested, tbl); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func parseColors(table *lua.LTable) (editor.Colors, error) {
	var colors editor.Colors
	err := parseStructTable(&colors, table)
	return colors, err
}

func parseProperties(table *lua.LTable) (editor.Properties, error) {
	var properties editor.Properties
	err := parseStructTable(&properties, table)
	return properties, err
}

func HandleSet(em *ExtensionManager, L *lua.LState) int {
	key := L.CheckString(1)

	switch key {
	case "colors":
		luaTable := L.CheckTable(2)
		colors, err := parseColors(luaTable)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
		em.M.Config.Colors = colors
		em.M.ReloadConfig()
		return 0

	case "properties":
		luaTable := L.CheckTable(2)
		properties, err := parseProperties(luaTable)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
		em.M.Config.Properties = properties
		em.M.ReloadConfig()
		return 0

	default:
		fieldName := snakeToCamel(key)
		cfgVal := reflect.ValueOf(&em.M.Config).Elem()

		propsVal := cfgVal.FieldByName("Properties")
		if propsVal.IsValid() {
			field := propsVal.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				value := L.CheckAny(2)

				switch field.Kind() {
				case reflect.String:
					if v, ok := value.(lua.LString); ok {
						field.SetString(string(v))
						em.M.ReloadConfig()
						return 0
					}

				case reflect.Bool:
					if v, ok := value.(lua.LBool); ok {
						field.SetBool(bool(v))
						em.M.ReloadConfig()
						return 0
					}

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if v, ok := value.(lua.LNumber); ok {
						field.SetInt(int64(v))
						em.M.ReloadConfig()
						return 0
					}

				case reflect.Float32, reflect.Float64:
					if v, ok := value.(lua.LNumber); ok {
						field.SetFloat(float64(v))
						em.M.ReloadConfig()
						return 0
					}
				}
			}
		}

		L.RaiseError("unknown config key: %s", key)
		return 0
	}
}

func GetConfigTable(em *ExtensionManager) *lua.LTable {
	config := em.L.NewTable()
	em.L.SetFuncs(config, map[string]lua.LGFunction{
		"set": func(L *lua.LState) int {
			return HandleSet(em, L)
		},
		"reload": func(L *lua.LState) int {
			em.M.ReloadConfig()
			return 0
		},
	})

	return config
}
