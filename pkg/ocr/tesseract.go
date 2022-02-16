package ocr

import (
	"fmt"
	"os"
	"time"

	"github.com/barklan/waysnip/pkg/system"
)

func Process(bytes []byte) (string, error) {
	f, err := os.CreateTemp("", "waysniptmp.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for image: %w", err)
	}
	defer os.Remove(f.Name())
	_, err = f.Write(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to write image to temporary file: %w", err)
	}

	// TODO "--psm 12" arguments should be under toggle (this is sparse text - good for code)
	out, ok, err := system.ExecIn(5*time.Second, "tesseract", "--psm", "12", "-l", "eng", f.Name(), "-")
	if !ok {
		return "", fmt.Errorf("ocr timeout")
	}
	if err != nil {
		return "", fmt.Errorf("failed to ocr image: %w", err)
	}
	return string(out), nil
}
