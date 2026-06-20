package logstuff

import "context"

type levelEnabler interface {
	Enabled(context.Context, LogLevel) bool
}
