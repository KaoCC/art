package main

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func ReadImage(path string) (image.Image, error) {

	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer imageFile.Close()

	rawData, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	return rawData, nil

}

func WriteImage(path string, imageData image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jpeg.Encode(file, imageData, nil); err != nil {
		return err
	}

	return nil
}

func WriteGif(path string, frames []image.Image) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	outGif := &gif.GIF{}
	for _, frame := range frames {
		bounds := frame.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.WebSafe)
		draw.Draw(palettedImage, palettedImage.Rect, frame, bounds.Min, draw.Over)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 0)
	}

	if err := gif.EncodeAll(file, outGif); err != nil {
		return err
	}

	return nil
}
