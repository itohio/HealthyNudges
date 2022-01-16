package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/HealthyNudges/pkg/config"
)

func makeSlider(min, max, step float64, format string) (binding.Float, *fyne.Container) {
	flt := binding.NewFloat()

	slider := widget.NewSlider(min, max)
	slider.Bind(flt)
	slider.Step = step

	sflt := binding.FloatToStringWithFormat(flt, format)
	label := widget.NewLabelWithData(sflt)
	return flt, container.NewBorder(nil, nil, nil, label, slider)
}

func (w *SettingsWindow) makeNudges() *container.TabItem {
	eName := widget.NewEntry()
	bWindow := widget.NewCheck("Show window", nil)
	bNotification := widget.NewCheck("Send notification", nil)
	bActive := widget.NewCheck("Nudge is active", func(b bool) {
		if !b {
			if bWindow.Checked {
				bWindow.SetChecked(false)
			}
			if bNotification.Checked {
				bNotification.SetChecked(false)
			}
		}
	})

	bWindow.OnChanged = func(b bool) {
		a, b := bNotification.Checked, bWindow.Checked
		if a || b {
			bActive.SetChecked(true)
		} else {
			bActive.SetChecked(false)
		}
	}
	bNotification.OnChanged = bWindow.OnChanged

	sType := widget.NewSelect(config.NudgeOptions, nil)
	sType.SetSelectedIndex(int(config.NudgeToType("")))

	fPeriods, sPeriods := makeSlider(1, 8, 1, "%2.0f min")
	fWork, sWork := makeSlider(0, 60, .5, "%2.1f min")
	fShort, sShort := makeSlider(3, 60, .5, "%2.1f min")
	fLong, sLong := makeSlider(3, 120, .5, "%2.1f min")

	fPeriods.Set(3)
	fWork.Set(45)
	fShort.Set(15)
	fLong.Set(20)

	tNudges := widget.NewListWithData(
		w.config.Nudges,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(2, widget.NewLabel(""), widget.NewLabel(""))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			eb, ok := di.(binding.Untyped)
			if !ok {
				return
			}
			ebv, err := eb.Get()
			if err != nil {
				return
			}
			nudge, ok := ebv.(*config.Nudge)
			if !ok {
				return
			}
			row := co.(*fyne.Container)
			t, ok := row.Objects[0].(*widget.Label)
			t.SetText(config.NudgeOptions[nudge.Type])
			n, ok := row.Objects[1].(*widget.Label)
			n.SetText(nudge.Name)
		},
	)

	nuSelected := -1
	tNudges.OnSelected = func(id widget.ListItemID) {
		nuSelected = id
		if id < 0 {
			return
		}
		nudge, ok := w.config.GetNudge(id)
		if !ok {
			return
		}
		eName.SetText(nudge.Name)
		bWindow.SetChecked(nudge.Window)
		bNotification.SetChecked(nudge.Notification)
		sType.SetSelectedIndex(int(nudge.Type))
		fWork.Set(nudge.Work)
		fShort.Set(nudge.ShortRest)
		fLong.Set(nudge.LongRest)
		fPeriods.Set(float64(nudge.WorkPeriods))
	}

	form := widget.NewForm()
	controls := container.NewHBox(
		widget.NewButton("Add", func() {
			w.addNudge(eName.Text, bWindow.Checked, bNotification.Checked, sType.Selected, fWork, fShort, fLong, fPeriods)
			form.Refresh()
		}),
		widget.NewButton("Update", func() {
			w.updateNudge(nuSelected, eName.Text, bWindow.Checked, bNotification.Checked, sType.Selected, fWork, fShort, fLong, fPeriods)
			tNudges.Refresh()
		}),
		widget.NewButton("Delete", func() {
			if nuSelected < 0 {
				return
			}
			items, err := w.config.Exceptions.Get()
			if err != nil {
				return
			}
			items = append(items[:nuSelected], items[nuSelected+1:]...)
			w.config.Exceptions.Set(items)
		}),
	)

	const ITEMS = 10
	items := [ITEMS]*widget.FormItem{
		widget.NewFormItem("", eName),
		widget.NewFormItem("", bWindow),
		widget.NewFormItem("", bNotification),
		widget.NewFormItem("", bActive),
		widget.NewFormItem("", sType),
		widget.NewFormItem("", sWork),
		widget.NewFormItem("", sShort),
		widget.NewFormItem("", sLong),
		widget.NewFormItem("", sPeriods),
		widget.NewFormItem("", controls),
	}
	for i, hint := range [ITEMS]string{
		"Nudge name",
		"Show a window",
		"Send a notification",
		"The nudge is active",
		"Type of the nudge",
		"Work duration",
		"Short rest duration",
		"Long rest duration",
		"Number of short rests before long rest",
		"",
	} {
		items[i].HintText = hint
	}

	form.Items = items[:]

	form.CancelText = "Defaults"
	form.SubmitText = "Save"

	form.OnCancel = func() {
		w.config.ReadExceptions()
		tNudges.Refresh()
	}
	form.OnSubmit = func() {
		w.config.WriteExceptions()
	}

	return container.NewTabItem("Nudges", container.NewBorder(nil, form, nil, nil, container.NewMax(tNudges)))
}

func (w *SettingsWindow) makeNudge(name string, window, notification bool, nudgeType string, fWork, fShort, fLong, fPeriods binding.Float) *config.Nudge {
	work, _ := fWork.Get()
	short, _ := fWork.Get()
	long, _ := fWork.Get()
	periods, _ := fWork.Get()

	e := &config.Nudge{
		Name:         name,
		Window:       window,
		Notification: notification,
		Type:         config.NudgeToType(nudgeType),
		Work:         work,
		ShortRest:    short,
		LongRest:     long,
		WorkPeriods:  int(periods),
	}

	return e
}

func (w *SettingsWindow) addNudge(name string, window, notification bool, nudgeType string, fWork, fShort, fLong, fPeriods binding.Float) {
	if name == "" {
		return
	}

	e := w.makeNudge(name, window, notification, nudgeType, fWork, fShort, fLong, fPeriods)
	w.config.Nudges.Append(e)
}

func (w *SettingsWindow) updateNudge(id int, name string, window, notification bool, nudgeType string, fWork, fShort, fLong, fPeriods binding.Float) {
	if id < 0 || name == "" {
		return
	}

	nudge, ok := w.config.GetNudge(id)
	if !ok {
		return
	}

	n := w.makeNudge(name, window, notification, nudgeType, fWork, fShort, fLong, fPeriods)
	n.Runtime = nudge.Runtime

	w.config.Nudges.SetValue(id, n)
}
