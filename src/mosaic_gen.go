package src

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/nfnt/resize"
)

type stats map[string]any 

func (s stats) String() string {
	var repr strings.Builder
	for k, v := range s {
		repr.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}
	return repr.String()
}

var Stats stats = make(stats)

func CreateMosaicPhoto(mainPhotoPath string, mosaicTilesDirPath string, tileProportion float64) (image.Image, stats, error) {	
	startTime := time.Now()
	mainPhoto, err := openImage(mainPhotoPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	tileSize := int(float64(mainPhoto.Bounds().Dx() + mainPhoto.Bounds().Dy()) * (tileProportion / 100))
	sdi := createSubdividedImage(mainPhoto, tileSize)
	tileImages := createTileImages(mosaicTilesDirPath, tileSize) 
	resImgBounds := image.Rect(0, 0, mainPhoto.Bounds().Dx() - mainPhoto.Bounds().Dx() % tileSize, mainPhoto.Bounds().Dy() - mainPhoto.Bounds().Dy() % tileSize)
	resultImage := image.NewRGBA(resImgBounds)
	var totalTilesUsed int
	for i, row := range sdi {
		for j, col := range row {
			tile := tileImages[nearestRGB(col, tileImages)]
			blitLocation := image.Rect(j * tileSize, i * tileSize, j * tileSize + tileSize, i * tileSize + tileSize)
			draw.Draw(resultImage, blitLocation, tile, image.Point{}, draw.Src)
			totalTilesUsed++
		}
	}
	Stats["Time to complete"] = time.Since(startTime)
	Stats["Total tiles used"] = totalTilesUsed
	return resultImage, Stats, nil
}

func openImage(imagePth string) (image.Image, error) {
	imageFile, err := os.Open(imagePth)
	if err != nil {
		return nil, err
	}
	imgData, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return imgData, nil
}

func imageResizeWrapper(img image.Image, newWidth uint, newHeight uint) image.Image{
	return resize.Resize(newWidth, newHeight, img, resize.Bilinear)
}

func rgbAvg(img image.Image) color.RGBA {
	var avgR, avgG, avgB, avgA uint32
	for i := 0; i < img.Bounds().Dy(); i++ {
		for j := 0; j < img.Bounds().Dx(); j++ {
			r, g, b, a := img.At(j, i).RGBA()
			avgR += r >> 8
			avgG += g >> 8
			avgB += b >> 8
			avgA += a >> 8
		}
	}
	totalPixels := img.Bounds().Dx() * img.Bounds().Dy()
	avgR /= uint32(totalPixels)
	avgG /= uint32(totalPixels)
	avgB /= uint32(totalPixels)
	avgA /= uint32(totalPixels)
	return color.RGBA{uint8(avgR), uint8(avgG), uint8(avgB), uint8(avgA)}
}

func writeImageToDisk(img image.Image, fileName string) {
	outFile, err := os.OpenFile(fileName, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664)
	if err != nil {
		panic(err)
	}
	jpeg.Encode(outFile, img, &jpeg.Options{jpeg.DefaultQuality})
}

func createSubdividedImage(img image.Image, tileSize int) [][]color.RGBA {
	sdi := make([][]color.RGBA, img.Bounds().Dy() / tileSize) 	
	for i := 0; i < img.Bounds().Dy() / tileSize; i++ {
		sdi[i] = make([]color.RGBA, img.Bounds().Dx() / tileSize)
		for j := 0; j < img.Bounds().Dx() / tileSize; j++ {
			subImgPoint := image.Point {j * tileSize, i * tileSize}
			subImg := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
			draw.Draw(subImg, subImg.Bounds(), img, subImgPoint, draw.Src)
			sdi[i][j] = rgbAvg(subImg)
		}
	}
	return sdi
}

func createTileImages(dirPath string, tileSize int) map[color.RGBA]image.Image {
	tiles := make(map[color.RGBA]image.Image)
	onFile := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !d.IsDir() {
			img, err := openImage(path)
			if err != nil {
				return nil
			}
			img = imageResizeWrapper(img, uint(tileSize), uint(tileSize))
			tiles[rgbAvg(img)] = img
		}
		return nil
	}
	err := filepath.WalkDir(dirPath, onFile)
	if err != nil {
		log.Fatal(err)
	}
	return tiles
}

func nearestRGB(rgb color.RGBA, tiles map[color.RGBA]image.Image) color.RGBA {
	var nearestRGBSoFar color.RGBA
	var minDistance float64
	for avgRGB, _ := range tiles {
		// hack to get random starting point
		nearestRGBSoFar = avgRGB
		minDistance = rgbDistance(rgb, nearestRGBSoFar)
		break
	} 
	for avgRGB, _ := range tiles {
		if temp := rgbDistance(rgb, avgRGB); temp < minDistance {
			minDistance = temp
			nearestRGBSoFar = avgRGB
		}
	}

	return nearestRGBSoFar
}

func rgbDistance(rgb1, rgb2 color.RGBA) float64 {
	// sqrt <- (x2−x1)^2+(y2−y1)^2+(z2−z1)^2
	x1, y1, z1 := float64(rgb1.R), float64(rgb1.G), float64(rgb1.B)
	x2, y2, z2 := float64(rgb2.R), float64(rgb2.G), float64(rgb2.B)
	return math.Sqrt(
		math.Pow((x2 - x1), 2) + math.Pow((y2 - y1), 2) + math.Pow((z2 - z1), 2),
	)
}


