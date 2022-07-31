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
	logs.Disable()
	enable := widget.NewCheckWithData("Enable logs", w.config.EnableLogs)
	return container.NewTabItem("Logs", container.NewBorder(enable, nil, nil, nil, logs))
}
