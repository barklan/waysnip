package pkg

import (
	"fmt"

	"github.com/barklan/waysnip/pkg/ocr"
	"github.com/barklan/waysnip/pkg/wlclip"
	"go.uber.org/zap"
)

func reportErr(lg *zap.Logger, err error) error {
	lg.Error("error will be copied to clipboard", zap.Error(err))
	return wlclip.ToClip(fmt.Sprintf("WAYSNIP ERROR: %s", err))
}

func Run(lg *zap.Logger) error {
	bb, err := wlclip.GetPNG()
	lg.Info("got png", zap.Int("nbytes", len(bb)))
	if err != nil {
		return reportErr(lg, err)
	}

	ready, err := ocr.PreProcess(bb)
	if err != nil {
		return reportErr(lg, err)
	}

	text, err := ocr.Process(ready)
	if err != nil {
		return reportErr(lg, err)
	}
	lg.Info("ocr successful", zap.String("text", text))

	return wlclip.ToClip(text)
}
