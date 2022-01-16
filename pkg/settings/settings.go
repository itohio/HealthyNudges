package settings

import (
	"github.com/itohio/HealthyNudges/pkg/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type SettingsWindow struct {
	fyne.Window
	app    fyne.App
	tabs   *container.AppTabs
	config *config.Config
}

func New(app fyne.App, config *config.Config) *SettingsWindow {
	ret := &SettingsWindow{
		app:    app,
		config: config,
	}
	ret.Window = app.NewWindow("Healthy Nudges Settings")
	ret.Window.Resize(fyne.NewSize(600, 800))

	// ret.Window.SetCloseIntercept(func() {
	// 	dialog.ShowConfirm(
	// 		"Do you want to exit?",
	// 		"Do you want to exit or hide?",
	// 		func(b bool) {
	// 			if b {
	// 				ret.Window.Close()
	// 			} else {
	// 			}
	// 		},
	// 		ret.Window,
	// 	)
	// })

	ret.tabs = container.NewAppTabs(
		ret.makeGeneral(),
		ret.makeNudges(),
		ret.makeExceptions(),
	)

	ret.Window.SetContent(ret.tabs)

	return ret
}
