package settings

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (w *SettingsWindow) makeLogs() *container.TabItem {
	logs := widget.NewMultiLineEntry()
	logs.Bind(w.bindLogs)
	logs.ActionItem = widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		w.bindLogs.Set("")
	})
	return container.NewTabItem("Logs", logs)
}
