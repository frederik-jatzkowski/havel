package statistics

import "context"

type Calculator interface {
	CalculateStatistics(ctx context.Context)
}
