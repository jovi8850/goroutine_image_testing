package imageprocessing

import (
	"image"
	"image/color"
	"os"
	"testing"
)

// createTestImage creates a simple test image for testing
func createTestImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 255 / width),
				G: uint8(y * 255 / height),
				B: 128,
				A: 255,
			})
		}
	}
	return img
}

func TestReadImage(t *testing.T) {
	// Create a temporary test image file
	tempFile := "test_image.jpg"
	defer os.Remove(tempFile)

	// Create and save test image
	testImg := createTestImage(100, 100)
	err := WriteImage(tempFile, testImg)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// Test reading the image
	img, err := ReadImage(tempFile)
	if err != nil {
		t.Errorf("ReadImage failed: %v", err)
	}

	if img == nil {
		t.Error("ReadImage returned nil image")
	}

	// Test reading non-existent file
	_, err = ReadImage("non_existent_file.jpg")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestWriteImage(t *testing.T) {
	tempFile := "test_output.jpg"
	defer os.Remove(tempFile)

	testImg := createTestImage(100, 100)

	err := WriteImage(tempFile, testImg)
	if err != nil {
		t.Errorf("WriteImage failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestGrayscale(t *testing.T) {
	testImg := createTestImage(100, 100)
	grayImg := Grayscale(testImg)

	if grayImg == nil {
		t.Error("Grayscale returned nil image")
	}

	// Verify it's a grayscale image by checking color model
	_, isGray := grayImg.(*image.Gray)
	if !isGray {
		t.Error("Grayscale did not return a grayscale image")
	}
}

func TestResize(t *testing.T) {
	testCases := []struct {
		name           string
		width          int
		height         int
		expectedMaxDim int
	}{
		{"Square image", 100, 100, 500},
		{"Landscape image", 800, 400, 500},
		{"Portrait image", 400, 800, 500},
		{"Small image", 50, 50, 500},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testImg := createTestImage(tc.width, tc.height)
			resizedImg, err := Resize(testImg)

			if err != nil {
				t.Errorf("Resize failed: %v", err)
			}

			if resizedImg == nil {
				t.Error("Resize returned nil image")
			}

			bounds := resizedImg.Bounds()
			newWidth := bounds.Dx()
			newHeight := bounds.Dy()

			// Check that aspect ratio is maintained (within rounding error)
			originalRatio := float64(tc.width) / float64(tc.height)
			newRatio := float64(newWidth) / float64(newHeight)
			ratioDiff := originalRatio - newRatio
			if ratioDiff < -0.1 || ratioDiff > 0.1 {
				t.Errorf("Aspect ratio not maintained: original %.2f, new %.2f", originalRatio, newRatio)
			}

			// Check that neither dimension exceeds max dimension
			if newWidth > 500 || newHeight > 500 {
				t.Errorf("Dimensions exceed max: %dx%d", newWidth, newHeight)
			}

			// Check that at least one dimension equals max dimension (for non-square images)
			if tc.width != tc.height && newWidth < 500 && newHeight < 500 {
				t.Errorf("No dimension reached max: %dx%d", newWidth, newHeight)
			}
		})
	}
}

func BenchmarkGrayscale(b *testing.B) {
	testImg := createTestImage(1000, 1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Grayscale(testImg)
	}
}

func BenchmarkResize(b *testing.B) {
	testImg := createTestImage(1000, 1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Resize(testImg)
	}
}

func BenchmarkReadWriteImage(b *testing.B) {
	tempFile := "benchmark_test.jpg"
	defer os.Remove(tempFile)

	testImg := createTestImage(500, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteImage(tempFile, testImg)
		ReadImage(tempFile)
	}
}
