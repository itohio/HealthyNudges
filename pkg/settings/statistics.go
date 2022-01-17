package settings

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (w *SettingsWindow) makeStatistics() *container.TabItem {
	const ITEMS = 0
	items := [ITEMS]*widget.FormItem{}

	for i, hint := range [ITEMS]string{} {
		items[i].HintText = hint
	}

	form := widget.NewForm(items[:]...)

	return container.NewTabItem("Statistics", form)
}
