package settings

import (
	"log"

	"github.com/itohio/HealthyNudges/pkg/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type SettingsWindow struct {
	fyne.Window
	app      fyne.App
	tabs     *container.AppTabs
	config   *config.Config
	bindLogs binding.String
}

func (s *SettingsWindow) Write(data []byte) (int, error) {
	str, _ := s.bindLogs.Get()
	s.bindLogs.Set(string(data) + str)
	return len(data), nil
}

func New(app fyne.App, config *config.Config) *SettingsWindow {
	ret := &SettingsWindow{
		app:      app,
		config:   config,
		bindLogs: binding.NewString(),
	}
	ret.Window = app.NewWindow("Healthy Nudges Settings")
	ret.Window.Resize(fyne.NewSize(600, 900))

	log.SetOutput(ret)

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
		ret.makeNudges(),
		ret.makeExceptions(),
		ret.makeGeneral(),
		//		ret.makeStatistics(),
		ret.makeLogs(),
	)

	ret.Window.SetContent(ret.tabs)
	return ret
}
