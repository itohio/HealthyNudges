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
	bFullscreen := widget.NewCheck("Fullscreen nudge window", nil)
	bAllScreens := widget.NewCheck("Cover all screens", nil)
	bWindow := widget.NewCheck("Show nudge window", nil)

	sAdvance := binding.FloatToStringWithFormat(w.config.ReminderAdvance, "%1.1f min")
	tReminderAdvance := widget.NewLabelWithData(sAdvance)

	bAutostart.Bind(w.config.AutoStart)
	bSystray.Bind(w.config.Systray)

	bReminder.Bind(w.config.ReminderNotification)
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

	items := []*widget.FormItem{
		widget.NewFormItem("", bAutostart),
		widget.NewFormItem("", bSystray),
		widget.NewFormItem("", bReminder),
		widget.NewFormItem("", container.NewBorder(nil, nil, nil, tReminderAdvance, sReminderAdvance)),
		widget.NewFormItem("", bWindow),
		widget.NewFormItem("", bFullscreen),
		widget.NewFormItem("", bAllScreens),
	}
	items[0].HintText = "Automatically starts when the user logs in."
	items[1].HintText = "Hide this window to Systray."
	items[2].HintText = "Sends system notifications before main period expires."
	items[3].HintText = "Send notifications the amount of minutes prior period expiry."
	items[4].HintText = "Displays a splash window with nudge information."
	items[5].HintText = "Show nudge window full screen."
	items[6].HintText = "Show nudge window on all screens."

	form := widget.NewForm(items...)
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
