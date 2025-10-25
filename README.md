# 🖼️ Go Image Processing Pipeline — Sequential vs Concurrent Comparison

## 📘 Assignment Overview

This project is based on **Amrit Singh’s Go Image Processing Pipeline (CODEHEIM)** example.  
It demonstrates how to process multiple images concurrently using Go’s **goroutines** and **channels**, and compares that to a sequential (non-concurrent) approach.

The pipeline performs these main operations:
1. **Read** image files from an input directory.  
2. **Resize** them (preserving aspect ratio).  
3. **Convert** them to grayscale.  
4. **Save** the processed images to an output folder.  
5. Measure performance for both **concurrent** and **sequential** modes.

---

## ⚙️ Summary of File Updates

### 🧩 `main.go` — Pipeline Control

| Feature | Original | Updated |
|----------|-----------|---------|
| **Error Handling** | Used `panic()` on file read/write errors. | Replaced with proper error returns and checks, allowing pipeline to continue on failure. |
| **Concurrency Toggle** | Always used goroutines for concurrent processing. | Added command-line flag `-concurrent` to toggle between **concurrent** and **sequential** modes. |
| **Sequential Mode** | Not implemented. | New `processSequentially()` function for non-concurrent runs. |
| **Performance Metrics** | No timing measurement. | Added `time.Since()` to measure total processing time. |
| **Status Reporting** | Printed only “Success!” or “Failed!”. | Reports per-file success/failure and totals at the end. |
| **CLI Flags** | None. | Added `flag.BoolVar()` and `flag.StringVar()` to support `-concurrent` and `-dir` options. |

#### 🔍 Example New CLI Usage:
```bash
# Concurrent processing (default)
.\imgproc.exe -concurrent=true -dir=images

# Sequential processing
.\imgproc.exe -concurrent=false -dir=images
```

---

### 🧠 `image_processing/image_processing.go` — Image Utilities

| Feature | Original | Updated |
|----------|-----------|---------|
| **Error Handling** | Panicked on file I/O errors. | Functions now return `error` values for better control. |
| **Resize Function** | Always forced images to 500x500 (distorted aspect ratio). | Now dynamically scales images to preserve aspect ratio (max dimension = 500px). |
| **Output Directory** | Assumed output folder exists. | Creates folders automatically with `os.MkdirAll()`. |
| **JPEG Encoding** | Used default quality settings. | Now saves JPEGs with quality level `90`. |

---

### 🧪 `image_processing/image_processing_test.go` — Unit and Benchmark Tests

Comprehensive test coverage was added for all helper functions using Go’s built-in `testing` framework.

#### ✅ Unit Tests
| Test Function | Purpose |
|----------------|----------|
| `TestReadImage` | Verifies successful file read and correct handling of missing files. |
| `TestWriteImage` | Ensures output files are properly created. |
| `TestGrayscale` | Confirms that grayscale conversion returns a valid `image.Gray` type. |
| `TestResize` | Validates resized image dimensions and ensures aspect ratio is preserved. |

#### ⚡ Benchmark Tests
| Benchmark | Purpose |
|------------|----------|
| `BenchmarkGrayscale` | Measures grayscale conversion performance. |
| `BenchmarkResize` | Measures resize performance. |
| `BenchmarkReadWriteImage` | Tests read/write throughput on disk operations. |

#### 🧰 How to Run Tests and Benchmarks
```powershell
# Run all unit tests
go test ./image_processing

# Run benchmarks (with memory and timing info)
go test -bench=. -benchmem ./image_processing
```

---

## 🧱 Project Structure

```
go_21_goroutines_pipeline/
│
├── images/
│   ├── image1.jpeg
│   ├── image2.jpeg
│   ├── image3.jpeg
│   ├── image4.jpeg
│   └── output/
│
├── image_processing/
│   ├── image_processing.go          # Core image functions (Read, Write, Resize, Grayscale)
│   ├── image_processing_test.go     # Unit and benchmark tests
│
├── main.go                          # Main pipeline with sequential & concurrent modes
├── go.mod, go.sum                   # Go module files
└── README.md                        # Project documentation
```

---

## 🧰 Building and Running the Program

### 1. Navigate to the project folder
```powershell
cd "[file path to folder]"
```

### 2. Install dependencies
```powershell
go mod tidy
```

### 3. Build the executable
```powershell
go build -o imgproc.exe
```

### 4. Run the program

#### Concurrent Mode
```powershell
.\imgproc.exe -concurrent=true -dir=images
```

#### Sequential Mode
```powershell
.\imgproc.exe -concurrent=false -dir=images
```

### 5. Check Results
Processed images are saved in:
```
images/output/
```

---

## ⏱️ Measuring Performance

To compare execution speeds:

```powershell
Measure-Command { .\imgproc.exe -concurrent=true }
Measure-Command { .\imgproc.exe -concurrent=false }
```

| Mode | Description | Expected Runtime |
|------|--------------|------------------|
| **Concurrent (-concurrent=true)** | Processes multiple images in parallel using goroutines | ⚡ Faster |
| **Sequential (-concurrent=false)** | Processes one image at a time | 🐢 Slower |

---

## 💡 Key Learnings

- Demonstrates **Go concurrency** via goroutines and channels.
- Uses **unit testing** and **benchmarking** for reliability and performance measurement.
- Implements **error handling** and **aspect ratio preservation** for image resizing.
- Compares **parallel vs sequential** processing for CPU and I/O-bound tasks.

---

## 🤖 GenAI Tools

This README and several code refactors were AI-assisted for clarity, test design, and documentation structure.  
All code was reviewed, debugged, and verified before submission. 
The repository contains the chats created in DeepSeek and ChatGPT for this assignment. 

---

## 🧾 Summary

This project now provides:
- A fully functional concurrent image pipeline with optional sequential mode.  
- Error-safe, testable, and modular image processing functions.  
- Automated test coverage and benchmarks.  
- Documented build, run, and performance comparison steps — meeting all assignment requirements.
