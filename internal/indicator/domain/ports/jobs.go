package ports

import (
	"context"
)

// IJobs Indicator Worker Runners Interface
type IJobs interface {
	CalculateByInterval(context.Context, string)
}
