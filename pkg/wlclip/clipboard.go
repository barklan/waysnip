package wlclip

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/barklan/waysnip/pkg/system"
	"go.uber.org/zap"
)

const (
	pngMime = "image/png"
	wlCopy  = "wl-copy"
	wlPaste = "wl-paste"
)

func ToClip(lg *zap.Logger, str string) error {
	cmd := exec.Command(wlCopy, "-n")
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to attach stdin pipe to wl-copy: %w", err)
	}
	if _, err = in.Write([]byte(str)); err != nil {
		return fmt.Errorf("failed to write to wl-copy: %w", err)
	}
	if err = in.Close(); err != nil {
		lg.Warn("failed to close input pipe to wl-copy", zap.Error(err))
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}
	return nil
}

func isImg() (bool, error) {
	out, ok, err := system.ExecIn(500*time.Millisecond, wlPaste, "-l")
	if !ok {
		return false, fmt.Errorf("wl-paste timeout (possibly empty clipboard?)")
	}
	if err != nil {
		return false, fmt.Errorf("wl-paste failed: %w, output: %s", err, out)
	}

	ok = bytes.Contains(out, []byte(pngMime))
	return ok, nil
}

func GetPNG() ([]byte, error) {
	ok, err := isImg()
	if err != nil {
		return nil, fmt.Errorf("failed to check clipboard: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("no image in clipboard")
	}

	out, ok, err := system.ExecIn(2*time.Second, wlPaste, "--type", pngMime)
	if !ok {
		return nil, fmt.Errorf("timeout getting image from clipboard")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get image from clipboard: %w", err)
	}
	return out, nil
}
