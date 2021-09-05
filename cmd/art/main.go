package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaocc/art/quadtree"
	// "github.com/pkg/profile"

	"flag"
)

func main() {

	/// CPU profiling
	// defer profile.Start().Stop()

	const defaultFile = "tux.jpg"
	const defaultSteps = 100
	fileName := defaultFile
	numSteps := 300
	isAnimated := false
	flag.StringVar(&fileName, "file", defaultFile, "The path to the target file")
	flag.IntVar(&numSteps, "step", defaultSteps, "Number of steps before stop")
	flag.BoolVar(&isAnimated, "animate", false, "Set to create an animated gif")
	helpFlag := flag.Bool("help", false, "Show usage")

	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	log.Printf("\n - Input File: %s\n - Steps: %d\n - Animated: %t\n", fileName, numSteps, isAnimated)

	inputImage, err := ReadImage(fileName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Build Tree\n")
	qtree := quadtree.QuadTree{}
	qtree.BuildTree(inputImage)

	log.Printf("Create Images\n")
	frames := qtree.CreateImages(numSteps, isAnimated)

	log.Printf("Output Images\n")

	outputFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	WriteImage(outputFileName+"_final.jpg", frames[len(frames)-1])

	if isAnimated {
		frames = append(frames, inputImage)
		WriteGif(outputFileName+"_animated.gif", frames)
	}

}
