package settings

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/HealthyNudges/pkg/config"
	"github.com/itohio/HealthyNudges/pkg/nudge"
)

func (w *SettingsWindow) makeExceptions() *container.TabItem {
	eName := widget.NewEntry()
	bActive := widget.NewCheck("Window must be active", nil)
	bExactMatch := widget.NewCheck("Exact match", nil)
	sType := widget.NewSelect(config.ExceptionOptions, nil)
	sHow := widget.NewSelect(config.HowOptions, nil)

	sType.SetSelectedIndex(int(config.ExceptionToType("")))
	sHow.SetSelectedIndex(int(config.HowToType("")))

	eName.Validator = func(val string) error {
		switch {
		case config.ExceptionType(sType.SelectedIndex()) == config.Times:
			return nudge.ValidateSchedule(val)
		case config.ExceptionType(sType.SelectedIndex()) == config.UserIdle:
			_, err := time.ParseDuration(val)
			return err
		case len(val) < 3:
			return fmt.Errorf("Name must be at least 3 symbols.")
		}
		return nil
	}

	form := widget.NewForm()
	tExceptions := widget.NewListWithData(
		w.config.Exceptions,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
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
			exception, ok := ebv.(*config.Exception)
			if !ok {
				return
			}
			row := co.(*fyne.Container)
			t, ok := row.Objects[0].(*widget.Label)
			t.SetText(config.ExceptionOptions[exception.Type])
			n, ok := row.Objects[1].(*widget.Label)
			n.SetText(exception.Name)
			h, ok := row.Objects[2].(*widget.Label)
			h.SetText(config.HowOptions[exception.How])
		},
	)

	exSelected := -1
	tExceptions.OnSelected = func(id widget.ListItemID) {
		exSelected = id
		if id < 0 {
			return
		}
		e, ok := w.config.GetException(id)
		if !ok {
			return
		}
		sType.SetSelectedIndex(int(e.Type))
		sHow.SetSelectedIndex(int(e.How))
		eName.SetText(e.Name)
		bActive.SetChecked(e.Active)
		bExactMatch.SetChecked(e.ExactMatch)
	}

	sType.OnChanged = func(s string) {
		form.Refresh()
	}

	controls := container.NewHBox(
		widget.NewButton("Add", func() {
			defer form.Refresh()
			if eName.Validator(eName.Text) != nil {
				return
			}
			w.addException(eName.Text, bActive.Checked, bExactMatch.Checked, sType.Selected, sHow.Selected)
		}),
		widget.NewButton("Update", func() {
			defer form.Refresh()
			if eName.Validator(eName.Text) != nil {
				return
			}
			w.updateException(exSelected, eName.Text, bActive.Checked, bExactMatch.Checked, sType.Selected, sHow.Selected)
			tExceptions.Refresh()
		}),
		widget.NewButton("Delete", func() {
			if exSelected < 0 {
				return
			}
			items, err := w.config.Exceptions.Get()
			if err != nil {
				return
			}
			items = append(items[:exSelected], items[exSelected+1:]...)
			w.config.Exceptions.Set(items)
		}),
	)

	const ITEMS = 6
	items := [ITEMS]*widget.FormItem{
		widget.NewFormItem("", eName),
		widget.NewFormItem("", bActive),
		widget.NewFormItem("", bExactMatch),
		widget.NewFormItem("", sType),
		widget.NewFormItem("", sHow),
		widget.NewFormItem("", controls),
	}

	for i, hint := range [ITEMS]string{
		"Exception details",
		"Window must be active(type = title)",
		"Title/Process exact match",
		"Type of the exception",
		"How to apply the exception",
		"",
	} {
		items[i].HintText = hint
	}

	form.Items = items[:]

	form.CancelText = "Defaults"
	form.SubmitText = "Save"

	form.OnCancel = func() {
		w.config.ReadExceptions()
		tExceptions.Refresh()
	}
	form.OnSubmit = func() {
		w.config.WriteExceptions()
	}

	return container.NewTabItem("Exceptions", container.NewBorder(nil, form, nil, nil, tExceptions))
}

func (w *SettingsWindow) makeException(name string, active, exactMatch bool, t, h string) *config.Exception {
	e := &config.Exception{
		Name:       name,
		Active:     active,
		ExactMatch: exactMatch,
		Type:       config.ExceptionToType(t),
		How:        config.HowToType(h),
	}

	return e
}

func (w *SettingsWindow) addException(name string, active, exactMatch bool, t, h string) {
	if name == "" {
		return
	}

	e := w.makeException(name, active, exactMatch, t, h)
	w.config.Exceptions.Append(e)
}

func (w *SettingsWindow) updateException(id int, name string, active, exactMatch bool, t, h string) {
	if id < 0 || name == "" {
		return
	}

	e := w.makeException(name, active, exactMatch, t, h)

	w.config.Exceptions.SetValue(id, e)
}
