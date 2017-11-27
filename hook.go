package zerolog

// Hook defines an interface to a log hook.
type Hook interface {
	// Run runs the hook with the event.
	Run(e *Event, level Level, message string)
}

// LevelHook applies a different hook for each level.
type LevelHook struct {
	NoLevelHook, DebugHook, InfoHook, WarnHook, ErrorHook, FatalHook, PanicHook Hook
}

// Run implements the Hook interface.
func (h LevelHook) Run(e *Event, level Level, message string) {
	switch level {
	case DebugLevel:
		if h.DebugHook != nil {
			h.DebugHook.Run(e, level, message)
		}
	case InfoLevel:
		if h.InfoHook != nil {
			h.InfoHook.Run(e, level, message)
		}
	case WarnLevel:
		if h.WarnHook != nil {
			h.WarnHook.Run(e, level, message)
		}
	case ErrorLevel:
		if h.ErrorHook != nil {
			h.ErrorHook.Run(e, level, message)
		}
	case FatalLevel:
		if h.FatalHook != nil {
			h.FatalHook.Run(e, level, message)
		}
	case PanicLevel:
		if h.PanicHook != nil {
			h.PanicHook.Run(e, level, message)
		}
	case NoLevel:
		if h.NoLevelHook != nil {
			h.NoLevelHook.Run(e, level, message)
		}
	}
}

// NewLevelHook returns a new LevelHook.
func NewLevelHook() LevelHook {
	return LevelHook{}
}
