package mosaicUtils

import (
	"Photo_Mosaic_Generation/src/image_types"
	// "fmt"
	"image"
	"github.com/nfnt/resize"
)

// mainPhotoPath is the path to the photo that the output image will be based off of.
// mosaicTilesDirPath path to a dir that holds photos that will be used to make up the mosaic tiles in
// the output image
func CreateMosaicPhoto(mainPhotoPath string, mosaicTilesDirPath string) (image.Image, error) {
	
	return nil, nil
}

// wrapper for resizing an image
func resizeWrapper(img image.Image, newWidth uint, newHeight uint) image.Image{
	return resize.Resize(newWidth, newHeight, img, resize.Bilinear)
}




