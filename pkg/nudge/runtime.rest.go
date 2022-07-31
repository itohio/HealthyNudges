package nudge

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/itohio/HealthyNudges/pkg/config"
)

func (r *nudgeRuntime) runRest(cfg *config.Config, nudge *config.Nudge, how config.ExceptionHow) {
	if r.stage == r.lastStage {
		return
	}

	if r.stage == Work {
		return
	}

	if r.stage == BreakReminder {
		r.NotifyReminder(cfg, nudge)
		return
	}

	r.restBreak(cfg, nudge)
}

func (r *nudgeRuntime) resetStage(nudge *config.Nudge) {
	r.counter++
	if r.counter >= nudge.WorkPeriods {
		r.counter = 0
	}
	r.stage = Work
	r.start = time.Now()
	r.stageTimer.Set(0)
}

func (r *nudgeRuntime) restBreak(cfg *config.Config, nudge *config.Nudge) {

	if !(r.stage == ShortBreak || r.stage == LongBreak) {
		return
	}

	if nudge.Notification {
		r.NotifyEvent(cfg, nudge)
	}

	if !nudge.Window {
		return
	}

	title := canvas.NewText(nudge.Name, theme.ForegroundColor())
	title.Alignment = fyne.TextAlignCenter
	title.TextSize = theme.TextHeadingSize()

	description := canvas.NewText(nudge.Description, theme.ForegroundColor())
	description.Alignment = fyne.TextAlignCenter
	description.TextSize = theme.TextSubHeadingSize()

	legend := canvas.NewText("", theme.PlaceHolderColor())
	legend.Alignment = fyne.TextAlignCenter
	legend.TextSize = theme.TextSize()

	im, err := loadImage("rest", fmt.Sprint(nudge.Name, " ", nudge.Description), cfg.GetMaxImageSize())
	if err != nil {
		log.Println("Failed to load an image: ", err)
		return
	}
	img := canvas.NewImageFromImage(im)
	img.FillMode = canvas.ImageFillOriginal

	pBar := widget.NewProgressBarWithData(r.stageTimer)
	pBar.Min = 0
	pBar.Max = r.stage.Duration(nudge)
	pBar.TextFormatter = func() string {
		val, _ := r.stageTimer.Get()
		if r.stage == Overtime {
			legend.Text = "Break is over!"
			legend.Refresh()
		}
		return fmt.Sprintf("%0.0f out of %0.0f min", val, pBar.Max)
	}

	button := widget.NewButton("Got it!", func() {
		r.Lock()
		defer r.Unlock()
		r.resetStage(nudge)
		if r.splash == nil {
			return
		}
		r.splash.Close()
	})
	content := container.NewVBox(
		title,
		description,
		container.NewPadded(container.NewCenter(img)),
		legend,
		pBar,
		button,
	)
	r.ShowEvent(cfg, nudge, content, func() {
		r.Lock()
		defer r.Unlock()
		r.resetStage(nudge)
	})
}
