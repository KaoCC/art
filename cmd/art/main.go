package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaocc/art/quadtree"
)

func main() {

	const defaultFile = "tux.jpg"
	const defaultSteps = 100

	var fileName string
	var numSteps uint
	var isAnimated bool
	var samplePeriod uint

	flag.StringVar(&fileName, "file", defaultFile, "The path to the target file.")
	flag.UintVar(&numSteps, "step", defaultSteps, "Number of steps before stop.")
	flag.BoolVar(&isAnimated, "animate", false, "Set to create an animated gif.")
	flag.UintVar(&samplePeriod, "period", 20, "The sample period. When animate flag is set, draw the result every n frame.")
	helpFlag := flag.Bool("help", false, "Show usage.")

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
	frames := qtree.CreateImages(numSteps, isAnimated, samplePeriod)

	log.Printf("Output Images\n")
	outputFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	err = WriteImage(outputFileName+"_final.jpg", frames[len(frames)-1])
	if err != nil {
		log.Fatal(err)
	}

	if isAnimated {
		log.Printf("Output Animated\n")
		frames = append(frames, inputImage)
		if err := WriteGif(outputFileName+"_animated.gif", frames); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Done\n")
}
