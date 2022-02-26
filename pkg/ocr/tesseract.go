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

	lang := "eng"
	if len(os.Args) >= 2 {
		lang = os.Args[1]
	}

	out, ok, err := system.ExecIn(5*time.Second, "tesseract", "-l", lang, f.Name(), "-")
	if !ok {
		return "", fmt.Errorf("ocr timeout")
	}
	if err != nil {
		return "", fmt.Errorf("failed to ocr image: %w", err)
	}
	return string(out), nil
}
