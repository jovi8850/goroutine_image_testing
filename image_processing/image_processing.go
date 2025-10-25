package imageprocessing

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func ReadImage(path string) (image.Image, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func WriteImage(path string, img image.Image) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}
	return nil
}

func Grayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel)
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

func Resize(img image.Image) (image.Image, error) {
	const maxDimension uint = 500
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())

	// Calculate new dimensions while maintaining aspect ratio
	var newWidth, newHeight uint
	if width > height {
		newWidth = maxDimension
		newHeight = uint(float64(height) * float64(maxDimension) / float64(width))
	} else {
		newHeight = maxDimension
		newWidth = uint(float64(width) * float64(maxDimension) / float64(height))
	}

	// Ensure minimum dimensions
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImg, nil
}
