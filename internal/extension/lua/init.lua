test = require("test")
test.run()

colors = {
	MainBackground = "#0E142E",
	MainForeground = "#dddddd",
	LineNumberBackground = "#090603",
	LineNumberForeground = "#dddddd",
	CursorColor = "#777777",
	CursorColorWrite = "#dddddd",
	PaletteBarBackground = "#3b3b3b",
	PaletteBarForeground = "#dddddd",
	PaletteInputBackground = "#090603",
	PaletteInputForeground = "#dddddd",
	InfoMsgBackground = "#090603",
	InfoMsgForeground = "#dddddd",
	WarnMsgBackground = "#090603",
	WarnMsgForeground = "#efc541",
	ErrorMsgBackground = "#090603",
	ErrorMsgForeground = "#ff5e56",
	WorkspaceBackground = "#3b3b3b",
	WorkspaceForeground = "#b9b9b9",
	WorkspaceBackgroundActive = "#dddddd",
	WorkspaceForegroundActive = "#3b3b3b",
	TabBackground = "#3b3b3b",
	TabForeground = "#b9b9b9",
	TabBackgroundActive = "#3b3b3b",
	TabForegroundActive = "#dddddd",
}

ms.config.set("properties", {
    GutterWidth = 4,
    TabSpaces = false,
})
ms.config.set("colors", colors)

-- ms.config.set()