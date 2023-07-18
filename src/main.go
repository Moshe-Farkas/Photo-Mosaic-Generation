/*
create new empty image and write it to disk:
img := image.NewRGBA(image.Rect(0, 0, 600, 200))
for i := 0; i < img.Rect.Dy()+1; i++ {
	for j := 0; j < img.Rect.Dx(); j++ {
		pixelVal := color.RGBA {}
		r, g, b := uint8(rand.Int()) % 255, uint8(rand.Int()) % 255, uint8(rand.Int()) % 255
		pixelVal.A = 255
		pixelVal.R = r
		pixelVal.G = g
		pixelVal.B = b
		img.Set(j, i, pixelVal)
	}
}
*/
package main

import (
	"fmt"
	"os"
	"errors"
	"image/jpeg"
	"Photo_Mosaic_Generation/src/mosaic_utils"
)
	
const EXPECTED_ARGS_LEN = 3

func main() {
	mainPhotoPath, mosaicTilesFolderPath, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	outputImage, err := mosaicUtils.CreateMosaicPhoto(mainPhotoPath, mosaicTilesFolderPath)
	
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	outputFile, err := os.OpenFile("out.jpg", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, outputImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		panic(err)
	}
}

func usage() string {
	return "Usage: mainPhoto.jpg folder/of/mosaic_tiles\n"
} 

func parseArgs() (string, string, error) {
	if len(os.Args) != EXPECTED_ARGS_LEN {
		return "", "", errors.New(usage())
	}
	mainPhotoPath := os.Args[1]
	if fileInfo, err := os.Stat(mainPhotoPath); err != nil {
		return "", "", errors.New("invalid main photo. Invalid path")
	} else if fileInfo.IsDir() {
		return "", "", errors.New("invalid main photo. Not a file")
	}	
	mosaicTilesFolderPath := os.Args[2]
	if fileInfo, err := os.Stat(mosaicTilesFolderPath); err != nil {
		return "", "", errors.New("invalid mosaic tiles folder. Invalid path")
	} else if !fileInfo.IsDir() {
		return "", "", errors.New("invalid mosaic tiles folder. Not a dir")
	}	
	return mainPhotoPath, mosaicTilesFolderPath, nil
}
