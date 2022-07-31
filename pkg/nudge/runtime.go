package nudge

import (
	"fmt"
	"log"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/itohio/HealthyNudges/pkg/config"
)

type nudgeStage int

const (
	Work nudgeStage = iota
	BreakReminder
	ShortBreak
	LongBreak
	Overtime
)

func (n nudgeStage) String() string {
	switch n {
	case Work:
		return "Work"
	case BreakReminder:
		return "Just before a break"
	case ShortBreak:
		return "Short break"
	case LongBreak:
		return "Long break"
	case Overtime:
		return "Break overtime"
	default:
		return "unknown"
	}
}

func (n nudgeStage) Duration(nudge *config.Nudge) float64 {
	switch n {
	case Work:
		return nudge.Work
	case ShortBreak:
		return nudge.ShortRest
	case LongBreak:
		fallthrough
	default:
		return nudge.LongRest
	}
}

type nudgeRuntime struct {
	sync.Mutex
	start      time.Time
	pauseStart time.Time
	notifiedAt time.Time
	windowAt   time.Time
	pauseDelay time.Duration
	counter    int
	paused     bool
	started    bool
	stage      nudgeStage
	lastStage  nudgeStage
	splash     fyne.Window
	stageTimer binding.Float
	stageName  binding.String
}

func newRuntime(nudge *config.Nudge) *nudgeRuntime {
	return &nudgeRuntime{
		start:      time.Now(),
		stageTimer: binding.NewFloat(),
		stageName:  binding.NewString(),
	}
}

func (r *nudgeRuntime) Run(cfg *config.Config, nudge *config.Nudge, how config.ExceptionHow) {
	if !(nudge.Window || nudge.Notification) {
		return
	}

	switch {
	case how == config.Stop && r.started:
		r.onStop(cfg, nudge)
	case how == config.Pause && !r.paused && r.started:
		r.onPause(cfg, nudge)
	case how == config.Run && !r.started:
		r.onStart(cfg, nudge)
	case how == config.Run && r.paused && r.started:
		r.onUnpause(cfg, nudge)
	}

	if r.paused || !r.started {
		return
	}

	r.calculateStage(cfg, nudge)

	if r.stage != r.lastStage {
		log.Println("Changed stage from ", r.lastStage, " to ", r.stage)
		r.stageName.Set(r.stage.String())
	}

	switch nudge.Type {
	case config.NudgeRest:
		r.runRest(cfg, nudge, how)
	case config.NudgeReminder:
		r.runReminder(cfg, nudge, how)
	}
}

func (r *nudgeRuntime) calculateStage(cfg *config.Config, nudge *config.Nudge) {
	t := time.Since(r.start) - r.pauseDelay
	stage := Work

	switch {
	case t < time.Duration(float64(time.Minute)*nudge.Work):
		rn, _ := cfg.ReminderNotification.Get()
		rb, _ := cfg.ReminderBeep.Get()
		if rn || rb {
			reminderAdvance, _ := cfg.ReminderAdvance.Get()
			if t >= time.Duration(float64(time.Minute)*(nudge.Work-reminderAdvance)) {
				stage = BreakReminder
			}
		}
		r.stageTimer.Set(t.Minutes())
	case t < time.Duration(float64(time.Minute)*(nudge.Work+nudge.ShortRest)) && r.counter < nudge.WorkPeriods:
		stage = ShortBreak
		r.stageTimer.Set(t.Minutes() - nudge.Work)
	case t < time.Duration(float64(time.Minute)*(nudge.Work+nudge.LongRest)) && r.counter >= nudge.WorkPeriods:
		r.stageTimer.Set(t.Minutes() - nudge.Work)
		stage = LongBreak
	default:
		r.stageTimer.Set(t.Minutes() - nudge.Work)
		stage = Overtime
	}

	r.lastStage, r.stage = r.stage, stage
}

func (r *nudgeRuntime) onStart(cfg *config.Config, nudge *config.Nudge) {
	r.start = time.Now()
	r.pauseDelay = 0
	r.counter = 0
	r.paused = false
	r.started = true
	log.Println("Nudge started: ", nudge.Name)
}

func (r *nudgeRuntime) onStop(cfg *config.Config, nudge *config.Nudge) {
	log.Println("Nudge stopped: ", nudge.Name)
	r.started = false
}

func (r *nudgeRuntime) onPause(cfg *config.Config, nudge *config.Nudge) {
	log.Println("Nudge paused: ", nudge.Name)
	r.pauseStart = time.Now()
	r.paused = true
}

func (r *nudgeRuntime) onUnpause(cfg *config.Config, nudge *config.Nudge) {
	log.Println("Nudge unpaused: ", nudge.Name)
	r.pauseDelay += time.Since(r.pauseStart)
	r.paused = false
}

// Remind a user about an event
func (r *nudgeRuntime) NotifyReminder(cfg *config.Config, nudge *config.Nudge) {
	adur, _ := cfg.ReminderAdvance.Get()
	if time.Since(r.notifiedAt).Minutes() < adur*2 {
		return
	}
	r.notifiedAt = time.Now()

	if ok, _ := cfg.ReminderNotification.Get(); ok {
		log.Println("Send Reminder Notification")
		advance, _ := cfg.ReminderAdvance.Get()
		fyne.CurrentApp().SendNotification(fyne.NewNotification(fmt.Sprintf("Healthy Nudge '%s' due in %0.1f min", nudge.Name, advance), nudge.Description))
	}
	if ok, _ := cfg.ReminderBeep.Get(); ok {
		// fyne.CurrentApp().SendNotification(fyne.NewNotification("Healthy Nudge Reminder: ", nudge.Description))
	}
}

// Notify the user about an event
func (r *nudgeRuntime) NotifyEvent(cfg *config.Config, nudge *config.Nudge) {
	adur, _ := cfg.ReminderAdvance.Get()
	if time.Since(r.notifiedAt).Minutes() < adur {
		return
	}
	r.notifiedAt = time.Now()

	if nudge.Notification {
		log.Println("Send Notification")
		fyne.CurrentApp().SendNotification(fyne.NewNotification(fmt.Sprintf("Healthy Nudge '%s'", nudge.Name), nudge.Description))
	}
}

func (r *nudgeRuntime) ShowEvent(cfg *config.Config, nudge *config.Nudge, content fyne.CanvasObject, onClosed func()) {
	if !nudge.Window {
		return
	}
	r.Lock()
	defer r.Unlock()

	if r.splash != nil {
		return
	}
	if time.Since(r.windowAt).Minutes() <= 1 {
		return
	}

	log.Println("Show Notify Window")

	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		_ = drv
		r.windowAt = time.Now()
		// r.splash = drv.CreateSplashWindow()
		r.splash = fyne.CurrentApp().NewWindow("Notification")
		r.splash.SetContent(content)
		// r.splash.Resize(fyne.NewSize(800, 800))
		r.splash.CenterOnScreen()
		r.splash.SetOnClosed(func() {
			if onClosed != nil {
				defer onClosed()
			}
			r.Lock()
			defer r.Unlock()
			r.splash = nil
			r.notifiedAt = time.Now()
		})
		// FIXME: Crashes
		//r.splash.RequestFocus()
		r.splash.Show()
	}
}
