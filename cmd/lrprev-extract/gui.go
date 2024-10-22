package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/schollz/progressbar/v3"
)

func runGUI2() {
	a := app.New()
	w := a.NewWindow("Lightroom Preview Extractor")

	inputDirEntry := widget.NewEntry()
	inputDirEntry.SetPlaceHolder("Input Directory")

	outputDirEntry := widget.NewEntry()
	outputDirEntry.SetPlaceHolder("Output Directory")

	lightroomDBEntry := widget.NewEntry()
	lightroomDBEntry.SetPlaceHolder("Lightroom Catalog Path")

	includeSizeCheck := widget.NewCheck("Include Image Size in Filename", nil)

	inputDirButton := widget.NewButton("Select Input Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				inputDirEntry.SetText(uri.Path())
			}
		}, w)
	})

	outputDirButton := widget.NewButton("Select Output Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				outputDirEntry.SetText(uri.Path())
			}
		}, w)
	})

	lightroomDBButton := widget.NewButton("Select Lightroom Catalog", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				lightroomDBEntry.SetText(uri.URI().Path())
			}
		}, w)
	})

	startButton := widget.NewButton("Start", func() {
		inputDir := inputDirEntry.Text
		outputDir := outputDirEntry.Text
		lightroomDB := lightroomDBEntry.Text
		includeSize := includeSizeCheck.Checked

		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		fileInfo, err := os.Stat(inputDir)
		if err != nil {
			log.Fatalf("Error accessing input path: %v", err)
		}

		if fileInfo.IsDir() {
			files, err := filepath.Glob(filepath.Join(inputDir, "**/*.lrprev"))
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
			err = processFile(inputDir, outputDir, lightroomDB, includeSize)
			if err != nil {
				log.Fatalf("Error processing file: %v", err)
			}
		}

		fmt.Println("Processing complete!")
	})

	w.SetContent(container.NewVBox(
		inputDirEntry,
		inputDirButton,
		outputDirEntry,
		outputDirButton,
		lightroomDBEntry,
		lightroomDBButton,
		includeSizeCheck,
		startButton,
	))

	w.ShowAndRun()
}

func runGUI() {
	a := app.New()
	w := a.NewWindow("Lightroom Preview Extractor")

	inputDirEntry := widget.NewEntry()
	inputDirEntry.SetPlaceHolder("Input Directory")

	outputDirEntry := widget.NewEntry()
	outputDirEntry.SetPlaceHolder("Output Directory")

	lightroomDBEntry := widget.NewEntry()
	lightroomDBEntry.SetPlaceHolder("Lightroom Catalog Path")

	includeSizeCheck := widget.NewCheck("Include Image Size in Filename", nil)

	startButton := widget.NewButton("Start", func() {
		inputDir := inputDirEntry.Text
		outputDir := outputDirEntry.Text
		lightroomDB := lightroomDBEntry.Text
		includeSize := includeSizeCheck.Checked

		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		fileInfo, err := os.Stat(inputDir)
		if err != nil {
			log.Fatalf("Error accessing input path: %v", err)
		}

		if fileInfo.IsDir() {
			files, err := filepath.Glob(filepath.Join(inputDir, "**/*.lrprev"))
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
			err = processFile(inputDir, outputDir, lightroomDB, includeSize)
			if err != nil {
				log.Fatalf("Error processing file: %v", err)
			}
		}

		fmt.Println("Processing complete!")
	})

	w.SetContent(container.NewVBox(
		inputDirEntry,
		outputDirEntry,
		lightroomDBEntry,
		includeSizeCheck,
		startButton,
	))

	w.ShowAndRun()
}