package nudge

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"github.com/itohio/HealthyNudges/pkg/config"
)

type Nudger struct {
	app    fyne.App
	config *config.Config
	ctx    context.Context
}

func New(app fyne.App, cfg *config.Config, ctx context.Context) *Nudger {
	ret := &Nudger{
		app:    app,
		config: cfg,
		ctx:    ctx,
	}

	return ret
}

func (s *Nudger) Start() {
	ticker := time.NewTicker(time.Second * 15)

	for {
		select {
		case <-ticker.C:
			s.Nudge()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Nudger) Nudge() {
	nudges, err := s.config.Nudges.Get()
	if err != nil {
		return
	}
	exceptions, err := s.config.Exceptions.Get()
	if err != nil {
		return
	}

	for _, nudge := range nudges {
		n, ok := nudge.(*config.Nudge)
		if !ok {
			continue
		}
		rt := s.runtime(n)
		if rt != nil {
			rt.Run(n, exceptions)
		}
	}
}

func (s *Nudger) runtime(nudge *config.Nudge) *nudgeRuntime {
	if nudge.Runtime == nil {
		s.makeRuntime(nudge)
	}

	rt, ok := nudge.Runtime.(*nudgeRuntime)
	if !ok {
		return nil
	}

	return rt
}

func (s *Nudger) makeRuntime(nudge *config.Nudge) {
	nudge.Runtime = newRuntime(nudge)
}
