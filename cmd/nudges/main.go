package main

import (
	"context"

	"fyne.io/fyne/v2/app"
	"github.com/itohio/HealthyNudges/pkg/config"
	"github.com/itohio/HealthyNudges/pkg/nudge"
	"github.com/itohio/HealthyNudges/pkg/settings"
)

func main() {
	myApp := app.NewWithID("itohio.healthy.nudges")
	cfg := config.New(myApp)
	wSettings := settings.New(myApp, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nudger := nudge.New(myApp, cfg, ctx)
	go nudger.Start()
	wSettings.ShowAndRun()
}
