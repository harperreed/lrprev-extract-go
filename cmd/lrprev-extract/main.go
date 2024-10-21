package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
	"lrprev-extract-go/internal/cli"
	"lrprev-extract-go/internal/extractor"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	inputDir := flag.String("d", "", "Path to your lightroom directory (.lrdata)")
	inputFile := flag.String("f", "", "Path to your file (.lrprev)")
	outputDirectory := flag.String("o", "", "Path to output directory")
	lightroomDB := flag.String("l", "", "Path to the lightroom catalog (.lrcat)")
	includeSize := flag.Bool("include-size", false, "Include image size information in the output file name")
	help := flag.Bool("help", false, "Show help information")
	gui := flag.Bool("gui", false, "Launch the GUI")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *gui {
		runGUI()
		return
	}

	if *inputDir == "" && *inputFile == "" {
		*inputDir = cli.PromptForInput("Enter the path to your lightroom directory (.lrdata) or file (.lrprev): ")
	}

	if *outputDirectory == "" {
		*outputDirectory = cli.PromptForInput("Enter the path to the output directory: ")
	}

	if *lightroomDB == "" {
		*lightroomDB = cli.PromptForInput("Enter the path to the lightroom catalog (.lrcat) [optional]: ")
	}

	if !*includeSize {
		*includeSize = cli.PromptForBool("Include image size information in the output file name? (y/n): ")
	}

	inputPath := *inputDir
	if *inputFile != "" {
		inputPath = *inputFile
	}

	err := os.MkdirAll(*outputDirectory, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatalf("Error accessing input path: %v", err)
	}

	if fileInfo.IsDir() {
		files, err := filepath.Glob(filepath.Join(inputPath, "**/*.lrprev"))
		if err != nil {
			log.Fatalf("Error finding .lrprev files: %v", err)
		}

		bar := progressbar.Default(int64(len(files)))

		for _, file := range files {
			err := processFile(file, *outputDirectory, *lightroomDB, *includeSize)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", file, err)
			}
			if err := bar.Add(1); err != nil {
				fmt.Printf("Error updating progress bar: %v\n", err)
			}
		}
	} else {
		err = processFile(inputPath, *outputDirectory, *lightroomDB, *includeSize)
		if err != nil {
			log.Fatalf("Error processing file: %v", err)
		}
	}

	fmt.Println("Processing complete!")
}

func processFile(filePath, outputDir, dbPath string, includeSize bool) error {
	fmt.Printf("Processing file: %s\n", filePath)
	return extractor.ExtractLargestJPEGFromLRPREV(filePath, outputDir, dbPath, includeSize)
}

func printHelp() {
	fmt.Println("lrprev-extract-go: Extract JPEG images from Lightroom preview files")
	fmt.Println("\nUsage:")
	fmt.Println("  lrprev-extract [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  lrprev-extract -d /path/to/lightroom/directory -o /path/to/output")
	fmt.Println("  lrprev-extract -f /path/to/file.lrprev -o /path/to/output -l /path/to/catalog.lrcat")
}

func runGUI() {
	a := app.New()
	w := a.NewWindow("LRPrev Extractor")

	inputDirEntry := widget.NewEntry()
	inputDirEntry.SetPlaceHolder("Path to your lightroom directory (.lrdata)")

	inputFileEntry := widget.NewEntry()
	inputFileEntry.SetPlaceHolder("Path to your file (.lrprev)")

	outputDirEntry := widget.NewEntry()
	outputDirEntry.SetPlaceHolder("Path to output directory")

	lightroomDBEntry := widget.NewEntry()
	lightroomDBEntry.SetPlaceHolder("Path to the lightroom catalog (.lrcat) [optional]")

	includeSizeCheck := widget.NewCheck("Include image size information in the output file name", nil)

	startButton := widget.NewButton("Start", func() {
		inputDir := inputDirEntry.Text
		inputFile := inputFileEntry.Text
		outputDir := outputDirEntry.Text
		lightroomDB := lightroomDBEntry.Text
		includeSize := includeSizeCheck.Checked

		inputPath := inputDir
		if inputFile != "" {
			inputPath = inputFile
		}

		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			log.Fatalf("Error accessing input path: %v", err)
		}

		if fileInfo.IsDir() {
			files, err := filepath.Glob(filepath.Join(inputPath, "**/*.lrprev"))
			if err != nil {
				log.Fatalf("Error finding .lrprev files: %v", err)
			}

			bar := progressbar.Default(int64(len(files)))

			for _, file := range files {
				err := processFile(file, outputDir, lightroomDB, includeSize)
				if err != nil {
					fmt.Printf("Error processing file %s: %v\n", file, err)
				}
				if err := bar.Add(1); err != nil {
					fmt.Printf("Error updating progress bar: %v\n", err)
				}
			}
		} else {
			err = processFile(inputPath, outputDir, lightroomDB, includeSize)
			if err != nil {
				log.Fatalf("Error processing file: %v", err)
			}
		}

		fmt.Println("Processing complete!")
	})

	w.SetContent(container.NewVBox(
		inputDirEntry,
		inputFileEntry,
		outputDirEntry,
		lightroomDBEntry,
		includeSizeCheck,
		startButton,
	))

	w.ShowAndRun()
}
