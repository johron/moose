package editor

import (
	"github.com/creasty/defaults"
	"github.com/gdamore/tcell/v3"
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
	MainBackground            string `default:"#090603"`
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
	WorkspaceBackgroundActive string `default:"#3b3b3b"`
	WorkspaceForegroundActive string `default:"#dddddd"`
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
