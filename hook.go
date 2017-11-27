package zerolog

// Hook defines an interface to a log hook.
type Hook interface {
	// Run runs the hook with the event.
	// hasLevel indicates whether this event is being logged
	// with a log level. It is false when Logger.Log() is used.
	Run(e *Event, level Level, hasLevel bool, message string)
}

// LevelHook applies a different hook for each level.
type LevelHook struct {
	NoLevelHook, DebugHook, InfoHook, WarnHook, ErrorHook, FatalHook, PanicHook Hook
}

// Run implements the Hook interface.
func (h LevelHook) Run(e *Event, level Level, hasLevel bool, message string) {
	if hasLevel {
		switch e.level {
		case DebugLevel:
			if h.DebugHook != nil {
				h.DebugHook.Run(e, level, hasLevel, message)
			}
		case InfoLevel:
			if h.InfoHook != nil {
				h.InfoHook.Run(e, level, hasLevel, message)
			}
		case WarnLevel:
			if h.WarnHook != nil {
				h.WarnHook.Run(e, level, hasLevel, message)
			}
		case ErrorLevel:
			if h.ErrorHook != nil {
				h.ErrorHook.Run(e, level, hasLevel, message)
			}
		case FatalLevel:
			if h.FatalHook != nil {
				h.FatalHook.Run(e, level, hasLevel, message)
			}
		case PanicLevel:
			if h.PanicHook != nil {
				h.PanicHook.Run(e, level, hasLevel, message)
			}
		}
	} else {
		if h.NoLevelHook != nil {
			h.NoLevelHook.Run(e, level, hasLevel, message)
		}
	}
}

// NewLevelHook returns a new LevelHook.
func NewLevelHook() LevelHook {
	return LevelHook{}
}
