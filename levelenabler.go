package logstuff

import "context"

// LevelEnabler is a minimal interface to gate logging behind a conditional.
type LevelEnabler interface {
	Enabled(context.Context, LogLevel) bool
}
