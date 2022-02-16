package ocr

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/anthonynsimon/bild/segment"
	"golang.org/x/image/draw"
)

func oversample(src image.Image) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X*2, src.Bounds().Max.Y*2))
	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst, nil
}

func toBW(src image.Image) image.Image {
	grayImg := image.NewGray(src.Bounds())
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			grayImg.Set(x, y, src.At(x, y))
		}
	}
	return grayImg
}

func threshold(img image.Image) image.Image {
	return segment.Threshold(img, 128)
}

func PreProcess(bb []byte) ([]byte, error) {
	src, err := png.Decode(bytes.NewBuffer(bb))
	if err != nil {
		return nil, fmt.Errorf("failed to decode png: %w", err)
	}

	img, err := oversample(src)
	if err != nil {
		return nil, fmt.Errorf("failed to oversample: %w", err)
	}

	img = toBW(img)

	img = threshold(img)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}

	// output, _ := os.Create("tmp.png")
	// defer output.Close()
	// _ = png.Encode(output, img)

	return buf.Bytes(), nil
}
