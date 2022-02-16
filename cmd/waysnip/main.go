package main

import (
	"github.com/barklan/waysnip/pkg"
	"github.com/barklan/waysnip/pkg/logging"
	"github.com/barklan/waysnip/pkg/system"
	"go.uber.org/zap"
)

func main() {
	go system.HandleSignals()
	internalEnv, _ := system.GetInternalEnv()

	lg := logging.New(internalEnv)
	defer func() {
		_ = lg.Sync()
	}()
	lg.Info("starting")
	defer lg.Info("main exited")

	if err := pkg.Run(lg); err != nil {
		lg.Panic("waysnip failed", zap.Error(err))
	}
}
