package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func makeSliderWithData(b binding.Float, min, max, step float64, format string) *fyne.Container {
	slider := widget.NewSlider(min, max)
	slider.Bind(b)
	slider.Step = step
	str := binding.FloatToStringWithFormat(b, format)
	label := widget.NewLabelWithData(str)
	return container.NewBorder(nil, nil, nil, label, slider)
}

func (w *SettingsWindow) makeGeneral() *container.TabItem {
	bAutostart := widget.NewCheckWithData("Run on start", w.config.AutoStart)
	bSystray := widget.NewCheckWithData("Hide to systray", w.config.Systray)
	bLogs := widget.NewCheckWithData("Enable logs", w.config.EnableLogs)
	sReminderAdvance := makeSliderWithData(w.config.ReminderAdvance, .5, .3, .5, "%1.1f min")
	bReminder := widget.NewCheckWithData("Reminder notifications", w.config.ReminderNotification)
	bReminderBeep := widget.NewCheckWithData("Reminder Beep", w.config.ReminderBeep)
	bFullscreen := widget.NewCheckWithData("Fullscreen nudge window", w.config.FullScreenWindow)
	bAllScreens := widget.NewCheckWithData("Cover all screens", w.config.AllScreens)
	sImageWidth := makeSliderWithData(w.config.MaxImageWidth, 100, 1024, 1, "%0.0f")
	sImageHeight := makeSliderWithData(w.config.MaxImageHeight, 100, 1024, 1, "%0.0f")
	bWindow := widget.NewCheckWithData("Show nudge window", w.config.ShowWindow)

	fWindow := bWindow.OnChanged
	bWindow.OnChanged = func(b bool) {
		if b {
			bFullscreen.Enable()
			bAllScreens.Enable()
			sImageWidth.Show()
			sImageHeight.Show()
		} else {
			bFullscreen.Disable()
			bAllScreens.Disable()
			sImageWidth.Hide()
			sImageHeight.Hide()
		}
		bFullscreen.Refresh()
		bAllScreens.Refresh()
		sImageWidth.Refresh()
		sImageHeight.Refresh()
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

	const ITEMS = 11
	items := [ITEMS]*widget.FormItem{
		widget.NewFormItem("", bAutostart),
		widget.NewFormItem("", bSystray),
		widget.NewFormItem("", bReminderBeep),
		widget.NewFormItem("", bReminder),
		widget.NewFormItem("", sReminderAdvance),
		widget.NewFormItem("", bWindow),
		widget.NewFormItem("", bFullscreen),
		widget.NewFormItem("", bAllScreens),
		widget.NewFormItem("", sImageWidth),
		widget.NewFormItem("", sImageHeight),
		widget.NewFormItem("", bLogs),
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
		"Maximum nudge image width.",
		"Maximum nudge image height.",
		"Enable collection of logs.",
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
