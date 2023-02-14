package ports

import (
	"context"
)

// IService Indicator domain service interface
type IService interface {
	CalculateByInterval(context.Context, string)
}
