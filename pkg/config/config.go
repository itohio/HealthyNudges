package config

import (
	"encoding/json"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const (
	REMINDER_NOTIFICATION    = "nudge.reminder"
	REMINDER_BEEP            = "nudge.reminder.beep"
	REMINDER_ADVANCE         = "nudge.reminder.advance"
	NUDGE_WINDOW_SHOW        = "nudge.window.show"
	NUDGE_WINDOW_FULLSCREEN  = "nudge.window.fullscreen"
	NUDGE_WINDOW_ALL_SCREENS = "nudge.window.allscreens"
	AUTO_START               = "app.autostart"
	HIDE_TO_SYSTRAY          = "app.systray"

	NUM_EXCEPTIONS = "except.num"
	FMT_EXCEPTION  = "except.%d"

	NUM_NUDGES = "nudge.num"
	FMT_NUDGE  = "nudge.%d"
)

type ExceptionType int
type ExceptionHow int
type NudgeType int

const (
	WindowTitle ExceptionType = iota
	Process
	Times
)
const (
	Run ExceptionHow = iota - 1
	Pause
	Stop
	Ignore
)
const (
	NudgeRest NudgeType = iota
	NudgeMeditate
	NudgeExcercise
	NudgeMeal
	NudgePomodoro
	NudgeReminder
)

var (
	ExceptionOptions = []string{"Window Title", "OS Process", "Times (cron format)"}
	HowOptions       = []string{"Pause all nudges", "Stop all nudges", "Ignore this exception"}
	NudgeOptions     = []string{"Rest regularly", "Meditate regularly", "Excercise regularly", "Have a healthy meal", "Pomodoro timer", "Reminder (cron format)"}
)

type Exception struct {
	Name       string        `json:"name"`
	Active     bool          `json:"active"`
	ExactMatch bool          `json:"exact_match"`
	Type       ExceptionType `json:"type"`
	How        ExceptionHow  `json:"how"`
}

type Nudge struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         NudgeType `json:"type"`
	Notification bool      `json:"notification"`
	Window       bool      `json:"window"`
	Schedule     string    `json:"schedule"`
	Work         float64   `json:"work_duration"`
	ShortRest    float64   `json:"short_rest_duration"`
	LongRest     float64   `json:"long_rest_duration"`
	WorkPeriods  int       `json:"periods_before_long_rest"`

	Runtime interface{}
}

type Config struct {
	app fyne.App

	// General config
	AutoStart binding.Bool
	Systray   binding.Bool

	// UI config
	ReminderNotification binding.Bool
	ReminderBeep         binding.Bool
	BeepSound            binding.String
	ShowWindow           binding.Bool
	FullScreenWindow     binding.Bool
	AllScreens           binding.Bool
	ReminderAdvance      binding.Float

	// Exceptions config
	Exceptions binding.UntypedList

	// Nudges config
	Nudges binding.UntypedList
}

func New(app fyne.App) *Config {
	ret := &Config{
		app:                  app,
		AutoStart:            binding.NewBool(),
		Systray:              binding.NewBool(),
		ReminderNotification: binding.NewBool(),
		ReminderBeep:         binding.NewBool(),
		ShowWindow:           binding.NewBool(),
		FullScreenWindow:     binding.NewBool(),
		AllScreens:           binding.NewBool(),
		ReminderAdvance:      binding.NewFloat(),
		Exceptions:           binding.NewUntypedList(),
		Nudges:               binding.NewUntypedList(),
	}
	ret.Read()

	return ret
}

func (c *Config) NotificationRollDur() time.Duration {
	val, _ := c.ReminderAdvance.Get()
	return time.Duration(float64(time.Minute) * val)
}

func (c *Config) Read() {
	c.ReadGeneral()
	c.ReadNudges()
	c.ReadExceptions()
}

func (c *Config) Write() {
	c.WriteGeneral()
	c.WriteNudges()
	c.WriteExceptions()
}

func (c *Config) ReadGeneral() {
	p := c.app.Preferences()
	c.AutoStart.Set(p.BoolWithFallback(AUTO_START, true))
	c.Systray.Set(p.BoolWithFallback(HIDE_TO_SYSTRAY, true))
	c.ReminderNotification.Set(p.BoolWithFallback(REMINDER_NOTIFICATION, true))
	c.ReminderBeep.Set(p.BoolWithFallback(REMINDER_BEEP, true))
	c.ShowWindow.Set(p.BoolWithFallback(NUDGE_WINDOW_SHOW, true))
	c.FullScreenWindow.Set(p.BoolWithFallback(NUDGE_WINDOW_FULLSCREEN, true))
	c.AllScreens.Set(p.BoolWithFallback(NUDGE_WINDOW_ALL_SCREENS, true))
	c.ReminderAdvance.Set(p.FloatWithFallback(REMINDER_ADVANCE, 1))
}

func (c *Config) WriteGeneral() {
	p := c.app.Preferences()
	bVal, _ := c.AutoStart.Get()
	p.SetBool(AUTO_START, bVal)
	bVal, _ = c.Systray.Get()
	p.SetBool(HIDE_TO_SYSTRAY, bVal)
	bVal, _ = c.ReminderNotification.Get()
	p.SetBool(REMINDER_NOTIFICATION, bVal)
	bVal, _ = c.ReminderBeep.Get()
	p.SetBool(REMINDER_BEEP, bVal)
	bVal, _ = c.ShowWindow.Get()
	p.SetBool(NUDGE_WINDOW_SHOW, bVal)
	bVal, _ = c.FullScreenWindow.Get()
	p.SetBool(NUDGE_WINDOW_FULLSCREEN, bVal)
	bVal, _ = c.AllScreens.Get()
	p.SetBool(NUDGE_WINDOW_ALL_SCREENS, bVal)
	fVal, _ := c.ReminderAdvance.Get()
	p.SetFloat(REMINDER_ADVANCE, fVal)
}

func (c *Config) ReadExceptions() {
	p := c.app.Preferences()
	num := p.IntWithFallback(NUM_EXCEPTIONS, 0)
	exceptions := make([]interface{}, 0, num)
	for i := 0; i < num; i++ {
		e := p.StringWithFallback(fmt.Sprintf(FMT_EXCEPTION, i), "{}")
		if e == "" {
			continue
		}
		exception := &Exception{}
		if json.Unmarshal([]byte(e), exception) == nil {
			exceptions = append(exceptions, exception)
		}
	}
	c.Exceptions.Set(exceptions)
}

func (c *Config) WriteExceptions() {
	p := c.app.Preferences()

	exceptions, _ := c.Exceptions.Get()
	n := 0
	for _, e := range exceptions {
		exception, ok := e.(*Exception)
		if !ok {
			continue
		}
		s, err := json.Marshal(exception)
		if err != nil {
			continue
		}

		p.SetString(fmt.Sprintf(FMT_EXCEPTION, n), string(s))

		n++
	}
	p.SetInt(NUM_EXCEPTIONS, n)
}

func (c *Config) ReadNudges() {
	p := c.app.Preferences()
	num := p.IntWithFallback(NUM_NUDGES, 0)
	nudges := make([]interface{}, 0, num)
	for i := 0; i < num; i++ {
		e := p.StringWithFallback(fmt.Sprintf(FMT_NUDGE, i), "{}")
		if e == "" {
			continue
		}
		nudge := &Nudge{}
		if json.Unmarshal([]byte(e), nudge) == nil {
			nudges = append(nudges, nudge)
		}
	}
	c.Nudges.Set(nudges)
}

func (c *Config) WriteNudges() {
	p := c.app.Preferences()

	nudges, _ := c.Nudges.Get()
	n := 0
	for _, e := range nudges {
		nudge, ok := e.(*Nudge)
		if !ok {
			continue
		}
		s, err := json.Marshal(nudge)
		if err != nil {
			continue
		}

		p.SetString(fmt.Sprintf(FMT_NUDGE, n), string(s))

		n++
	}
	p.SetInt(NUM_NUDGES, n)
}

func (c *Config) GetException(id int) (*Exception, bool) {
	e, err := c.Exceptions.GetItem(id)
	if err != nil {
		return nil, false
	}
	ev, ok := e.(binding.Untyped)
	if !ok {
		return nil, false
	}
	evb, err := ev.Get()
	if err != nil {
		return nil, false
	}
	exception, ok := evb.(*Exception)
	return exception, ok
}

func (c *Config) GetNudge(id int) (*Nudge, bool) {
	e, err := c.Nudges.GetItem(id)
	if err != nil {
		return nil, false
	}
	ev, ok := e.(binding.Untyped)
	if !ok {
		return nil, false
	}
	evb, err := ev.Get()
	if err != nil {
		return nil, false
	}
	nudge, ok := evb.(*Nudge)
	return nudge, ok
}

func OptionsToIdx(options []string, t string, def int) int {
	for i, s := range options {
		if s == t {
			return i
		}
	}
	return def
}

func ExceptionToType(t string) ExceptionType {
	return ExceptionType(OptionsToIdx(ExceptionOptions, t, int(WindowTitle)))
}

func HowToType(t string) ExceptionHow {
	return ExceptionHow(OptionsToIdx(HowOptions, t, int(Pause)))
}

func NudgeToType(t string) NudgeType {
	return NudgeType(OptionsToIdx(NudgeOptions, t, int(NudgeRest)))
}
