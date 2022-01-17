package nudge

import (
	"github.com/itohio/HealthyNudges/pkg/config"
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
		case config.WindowTitle:
		case config.Times:
			how = checkTimesException(exception)
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
