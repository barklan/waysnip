package system

import (
	"context"
	"os/exec"
	"time"
)

func ExecIn(deadline time.Duration, cmd string, args ...string) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	c := exec.CommandContext(ctx, cmd, args...)

	out, err := c.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return out, false, nil
	}

	return out, true, err
}
