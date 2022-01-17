package nudge

import (
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/itohio/HealthyNudges/pkg/config"
)

func (r *nudgeRuntime) runReminder(cfg *config.Config, nudge *config.Nudge, how config.ExceptionHow) {
	adur, _ := cfg.ReminderAdvance.Get()
	advance := time.Duration(float64(time.Minute) * adur)
	reminder := MatchSchedule(nudge.Schedule, advance)
	event := MatchSchedule(nudge.Schedule, 0)

	if reminder && !event {
		r.NotifyReminder(cfg, nudge)
	}

	if event {
		r.NotifyEvent(cfg, nudge)
		r.ShowEvent(cfg, nudge, widget.NewLabel("Hello World"))
	}
}
