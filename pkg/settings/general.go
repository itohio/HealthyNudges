package settings

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func (w *SettingsWindow) makeGeneral() *container.TabItem {
	bAutostart := widget.NewCheck("Run on start", nil)
	bSystray := widget.NewCheck("Hide to systray", nil)
	sReminderAdvance := widget.NewSlider(.5, 3)
	sReminderAdvance.Bind(w.config.ReminderAdvance)
	sReminderAdvance.Step = 0.5
	bReminder := widget.NewCheck("Reminder notifications", nil)
	bReminderBeep := widget.NewCheck("Reminder Beep", nil)
	bFullscreen := widget.NewCheck("Fullscreen nudge window", nil)
	bAllScreens := widget.NewCheck("Cover all screens", nil)
	bWindow := widget.NewCheck("Show nudge window", nil)

	sAdvance := binding.FloatToStringWithFormat(w.config.ReminderAdvance, "%1.1f min")
	tReminderAdvance := widget.NewLabelWithData(sAdvance)

	bAutostart.Bind(w.config.AutoStart)
	bSystray.Bind(w.config.Systray)

	bReminder.Bind(w.config.ReminderNotification)
	bReminderBeep.Bind(w.config.ReminderBeep)
	bFullscreen.Bind(w.config.FullScreenWindow)
	bAllScreens.Bind(w.config.AllScreens)
	bWindow.Bind(w.config.ShowWindow)

	fWindow := bWindow.OnChanged
	bWindow.OnChanged = func(b bool) {
		if b {
			bFullscreen.Enable()
		} else {
			bFullscreen.Disable()
		}
		bFullscreen.Refresh()
		fWindow(b)
	}
	fReminder := bReminder.OnChanged
	bReminder.OnChanged = func(b bool) {
		if b {
			sReminderAdvance.Show()
		} else {
			sReminderAdvance.Hide()
		}
		sReminderAdvance.Refresh()
		fReminder(b)
	}

	bAutostart.Disable()
	bSystray.Disable()

	const ITEMS = 8
	items := [ITEMS]*widget.FormItem{
		widget.NewFormItem("", bAutostart),
		widget.NewFormItem("", bSystray),
		widget.NewFormItem("", bReminderBeep),
		widget.NewFormItem("", bReminder),
		widget.NewFormItem("", container.NewBorder(nil, nil, nil, tReminderAdvance, sReminderAdvance)),
		widget.NewFormItem("", bWindow),
		widget.NewFormItem("", bFullscreen),
		widget.NewFormItem("", bAllScreens),
	}

	for i, hint := range [ITEMS]string{
		"Automatically starts when the user logs in.",
		"Hide this window to Systray.",
		"Remind about a nudge using sound.",
		"Remind about a nudge using notifications.",
		"Remind about a nudge a set amount of minutes prior.",
		"Displays a splash window with nudge information.",
		"Show nudge window full screen.",
		"Show nudge window on all screens.",
	} {
		items[i].HintText = hint
	}

	form := widget.NewForm(items[:]...)
	form.CancelText = "Defaults"
	form.SubmitText = "Save"

	form.OnCancel = func() {
		w.config.ReadGeneral()
		form.Refresh()
	}
	form.OnSubmit = func() {
		w.config.WriteGeneral()
	}

	return container.NewTabItem("General", form)
}
