package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaocc/art/quadtree"
	// "github.com/pkg/profile"

	"flag"
)

func main() {

	// CPU profiling by default
	// defer profile.Start().Stop()

	const defaultFile = "tux.jpg"
	fileName := defaultFile
	numSteps := 300
	isAnimated := false
	flag.StringVar(&fileName, "file", defaultFile, "a string var")
	flag.IntVar(&numSteps, "step", 100, "number of steps before stop")
	flag.BoolVar(&isAnimated, "animate", false, "set to create an animated gif")
	helpFlag := flag.Bool("help", false, "show usage")

	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf(" - Input File: %s\n - Steps: %d\n - Animated: %t\n", fileName, numSteps, isAnimated)

	inputImage, err := ReadImage(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Build Tree\n")
	qtree := quadtree.QuadTree{}
	qtree.BuildTree(inputImage)
	// qtree.Traversal()

	fmt.Printf("Create Images\n")
	frames := qtree.CreateImages(numSteps, isAnimated)

	fmt.Printf("Output Images\n")

	outputFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	WriteImage(outputFileName+"_final.jpg", frames[len(frames)-1])

	if isAnimated {
		// append original image
		frames = append(frames, inputImage)
		WriteGif(outputFileName+"_animated.gif", frames)
	}

}
