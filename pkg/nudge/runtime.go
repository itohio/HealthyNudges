package nudge

import (
	"time"

	"github.com/itohio/HealthyNudges/pkg/config"
)

type nudgeRuntime struct {
	start      time.Time
	pauseStart time.Time
	pauseDelay time.Duration
	paused     bool
}

func newRuntime(nudge *config.Nudge) *nudgeRuntime {
	return &nudgeRuntime{
		start: time.Now(),
	}
}

func (r *nudgeRuntime) Run(nudge *config.Nudge, exceptions []interface{}) {
	switch nudge.Type {
	case config.NudgeRest:
		fallthrough
	case config.NudgeMeditate:
		r.runRest(nudge, exceptions)
	}
}

func (r *nudgeRuntime) checkExceptions(nudge *config.Nudge, exceptions []interface{}) config.ExceptionHow {
	return config.Pause
}
