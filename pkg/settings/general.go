package settings

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func (w *SettingsWindow) makeGeneral() *container.TabItem {
	bAutostart := widget.NewCheck("Run on start", nil)
	bSystray := widget.NewCheck("Hide to systray", nil)
	sNotificationRoll := widget.NewSlider(.5, 3)
	sNotificationRoll.Bind(w.config.NotificationRoll)
	sNotificationRoll.Step = 0.5
	bNotifications := widget.NewCheck("Send notifications", nil)
	bFullscreen := widget.NewCheck("Fullscreen nudge window", nil)
	bAllScreens := widget.NewCheck("Cover all screens", nil)
	bWindow := widget.NewCheck("Show nudge window", nil)

	sRoll := binding.FloatToStringWithFormat(w.config.NotificationRoll, "%1.1f min")
	tNotificationRoll := widget.NewLabelWithData(sRoll)

	bAutostart.Bind(w.config.AutoStart)
	bSystray.Bind(w.config.Systray)

	bNotifications.Bind(w.config.SendNotifications)
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
	fNotifications := bNotifications.OnChanged
	bNotifications.OnChanged = func(b bool) {
		if b {
			sNotificationRoll.Show()
		} else {
			sNotificationRoll.Hide()
		}
		sNotificationRoll.Refresh()
		fNotifications(b)
	}

	bAutostart.Disable()
	bSystray.Disable()

	items := []*widget.FormItem{
		widget.NewFormItem("", bAutostart),
		widget.NewFormItem("", bSystray),
		widget.NewFormItem("", bNotifications),
		widget.NewFormItem("", container.NewBorder(nil, nil, nil, tNotificationRoll, sNotificationRoll)),
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
