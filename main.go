package main

import (
	"PHOTO_MOSAIC_GENERATION/src"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"os"
)

const EXPECTED_ARGS_LEN = 3

func main() {
	mainPhotoPath, mosaicTilesFolderPath, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	outputImage, stats, err := src.CreateMosaicPhoto(mainPhotoPath, mosaicTilesFolderPath, 0.3)
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
		log.Fatal(err.Error())
	}
	fmt.Println(stats.String())
	fmt.Println("finished")
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

