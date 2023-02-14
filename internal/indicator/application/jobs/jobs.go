package jobs

import (
	"context"
	"main/config"
	"main/internal/indicator/domain/ports"
	"main/pkg/logger"
	"main/pkg/utils/common"
	"time"
)

// jobRunner Indicator Worker struct
type jobRunner struct {
	cfg    *config.Config
	logger logger.Logger
	srv    ports.IService
}

// NewJobRunner Indicator worker constructor
func NewJobRunner(cfg *config.Config, logger logger.Logger, srv ports.IService) ports.IJobs {
	return &jobRunner{
		cfg:    cfg,
		logger: logger,
		srv:    srv,
	}
}

// CalculateByInterval calculate some works
func (w *jobRunner) CalculateByInterval(ctx context.Context, intervals string) {

	w.logger.Infof("Domain: %s, Mode: %s, Func: %s, Status: %s \n", "Indicator", "Worker", "Indicator", "Init")

	var counter = 0
	for range time.Tick(time.Second * common.StringToDuration(intervals)) {
		w.srv.CalculateByInterval(ctx, intervals)
		w.logger.Infof("CalculateByInterval Run Count : %d", counter)
		counter += 1
	}
}
