package nudge

import (
	"strings"
	"time"

	"github.com/itohio/HealthyNudges/pkg/config"
	"github.com/lextoumbourou/idle"
	"github.com/shirou/gopsutil/process"
)

func (s *Nudger) checkExceptions(exceptions []interface{}) config.ExceptionHow {
	for _, ei := range exceptions {
		exception, ok := ei.(*config.Exception)
		if !ok {
			continue
		}
		if exception.How == config.Ignore {
			continue
		}

		how := config.Ignore
		switch exception.Type {
		case config.Process:
			how = checkProcessException(exception)
		case config.WindowTitle:
			how = checkWindowTitleException(exception)
		case config.Times:
			how = checkTimesException(exception)
		case config.UserIdle:
			how = checkUserIdleException(exception)
		}
		if how != config.Ignore {
			return how
		}
	}

	return config.Run
}

func checkTimesException(e *config.Exception) config.ExceptionHow {
	if MatchSchedule(e.Name, 0) {
		return e.How
	}
	return config.Ignore
}

func checkProcessException(e *config.Exception) config.ExceptionHow {
	processes, err := process.Processes()
	if err != nil {
		return config.Ignore
	}
	eName := strings.ToLower(e.Name)
	for _, p := range processes {
		if ok, err := p.IsRunning(); !ok || err != nil {
			continue
		}
		pName, err := p.Name()
		if err != nil {
			continue
		}
		fg, err := p.Foreground()
		if err != nil {
			continue
		}
		if e.Active && !fg {
			continue
		}
		if e.ExactMatch {
			if pName == e.Name {
				return e.How
			}
		} else {
			if strings.Contains(strings.ToLower(pName), eName) {
				return e.How
			}
		}
	}
	return config.Ignore
}

func checkUserIdleException(e *config.Exception) config.ExceptionHow {
	idleTime, err := time.ParseDuration(e.Name)
	if err != nil {
		return config.Ignore
	}
	userIdleTime, err := idle.Get()
	if err != nil {
		return config.Ignore
	}
	if userIdleTime > idleTime {
		return e.How
	}
	return config.Ignore
}

func checkWindowTitleException(e *config.Exception) config.ExceptionHow {
	activeTitle := FindActiveWindowTitle()
	if e.ExactMatch {
		if activeTitle == e.Name {
			return e.How
		}
	} else {
		if strings.Contains(strings.ToLower(activeTitle), strings.ToLower(e.Name)) {
			return e.How
		}
	}
	return config.Ignore
}
