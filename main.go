package main

import (
	"flag"
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"strings"
	"time"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Error     error
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		for _, p := range paths {
			job := Job{
				InputPath: p,
				OutPath:   strings.Replace(p, "images/", "images/output/", 1),
			}
			job.Image, job.Error = imageprocessing.ReadImage(p)
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			if job.Error != nil {
				out <- job
				continue
			}
			job.Image, job.Error = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			if job.Error != nil {
				out <- job
				continue
			}
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			if job.Error != nil {
				out <- job
				continue
			}
			job.Error = imageprocessing.WriteImage(job.OutPath, job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func processConcurrently(imagePaths []string) time.Duration {
	start := time.Now()

	channel1 := loadImage(imagePaths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	successCount := 0
	failureCount := 0
	for job := range writeResults {
		if job.Error != nil {
			fmt.Printf("Failed to process %s: %v\n", job.InputPath, job.Error)
			failureCount++
		} else {
			fmt.Printf("Successfully processed %s -> %s\n", job.InputPath, job.OutPath)
			successCount++
		}
	}

	fmt.Printf("Concurrent processing completed: %d successes, %d failures\n", successCount, failureCount)
	return time.Since(start)
}

func processSequentially(imagePaths []string) time.Duration {
	start := time.Now()

	successCount := 0
	failureCount := 0

	for _, path := range imagePaths {
		outPath := strings.Replace(path, "images/", "images/output/", 1)

		// Load image
		img, err := imageprocessing.ReadImage(path)
		if err != nil {
			fmt.Printf("Failed to load %s: %v\n", path, err)
			failureCount++
			continue
		}

		// Resize image
		img, err = imageprocessing.Resize(img)
		if err != nil {
			fmt.Printf("Failed to resize %s: %v\n", path, err)
			failureCount++
			continue
		}

		// Convert to grayscale
		img = imageprocessing.Grayscale(img)

		// Save image
		err = imageprocessing.WriteImage(outPath, img)
		if err != nil {
			fmt.Printf("Failed to save %s: %v\n", outPath, err)
			failureCount++
		} else {
			fmt.Printf("Successfully processed %s -> %s\n", path, outPath)
			successCount++
		}
	}

	fmt.Printf("Sequential processing completed: %d successes, %d failures\n", successCount, failureCount)
	return time.Since(start)
}

func main() {
	var useConcurrent bool
	var imageDir string

	flag.BoolVar(&useConcurrent, "concurrent", true, "Use concurrent processing (true) or sequential processing (false)")
	flag.StringVar(&imageDir, "dir", "images", "Directory containing images")
	flag.Parse()

	imagePaths := []string{
		fmt.Sprintf("%s/image1.jpeg", imageDir),
		fmt.Sprintf("%s/image2.jpeg", imageDir),
		fmt.Sprintf("%s/image3.jpeg", imageDir),
		fmt.Sprintf("%s/image4.jpeg", imageDir),
	}

	fmt.Printf("Processing %d images...\n", len(imagePaths))
	fmt.Printf("Mode: %s\n", map[bool]string{true: "CONCURRENT", false: "SEQUENTIAL"}[useConcurrent])

	var duration time.Duration
	if useConcurrent {
		duration = processConcurrently(imagePaths)
	} else {
		duration = processSequentially(imagePaths)
	}

	fmt.Printf("Total processing time: %v\n", duration)
}
