package src

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"
)

func TestSubdividedImageLength(t *testing.T) {
	mockImage := image.NewRGBA(image.Rect(0, 0, 550, 650))
	tileSize := 100
	wantWidth := 5
	wantHeight := 6
	sdi := createSubdividedImage(mockImage, tileSize)
	gotHeight := len(sdi)
	gotWidth := len(sdi[gotHeight-1])
	if gotHeight != wantHeight {
		t.Fatalf("incorrect height")
	}
	if gotWidth != wantWidth {
		t.Fatalf("incorrect width")
	}
}

func TestSubdividedImageProperRGBAvg(t *testing.T) {
	tileSize := 54
	mockImage := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	for i := 0; i < tileSize; i++ {
		for j := 0; j < tileSize; j++ {
			mockImage.Set(j, i, color.RGBA{255, 0, 0, 255})
		}
	}
	want := color.RGBA{255, 0, 0, 255}
	got := rgbAvg(mockImage)
	if want != got {
		t.Fatalf("wrong avg")
	}
}

func TestCorrectRGBAvgPlacment(t *testing.T) {
	tileSize := 100
	blueMockImg := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(
		blueMockImg,
		blueMockImg.Bounds(),
		&image.Uniform{color.RGBA{0, 0, 255, 255}},
		image.Point{},
		draw.Src,
	)
	redMockImg := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(
		redMockImg,
		redMockImg.Bounds(),
		&image.Uniform{color.RGBA{255, 0, 0, 255}},
		image.Point{},
		draw.Src,
	)
	greenMockImg := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(
		greenMockImg,
		greenMockImg.Bounds(),
		&image.Uniform{color.RGBA{0, 255, 0, 255}},
		image.Point{},
		draw.Src,
	)
	blackMockImg := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(
		blackMockImg,
		blackMockImg.Bounds(),
		&image.Uniform{color.RGBA{0, 0, 0, 255}},
		image.Point{},
		draw.Src,
	)
	completeMockImg := image.NewRGBA(image.Rect(0, 0, tileSize*2, tileSize*2))
	// [blue][red]
	// [green][black]
	draw.Draw(
		completeMockImg,
		image.Rect(0, 0, tileSize, tileSize),
		blueMockImg,
		image.Point{},
		draw.Src,
	)
	draw.Draw(
		completeMockImg,
		image.Rect(tileSize, 0, tileSize*2, tileSize),
		redMockImg,
		image.Point{},
		draw.Src,
	)
	draw.Draw(
		completeMockImg,
		image.Rect(0, tileSize, tileSize*2, tileSize*2),
		greenMockImg,
		image.Point{},
		draw.Src,
	)
	draw.Draw(
		completeMockImg,
		image.Rect(tileSize, tileSize, tileSize*2, tileSize*2),
		blackMockImg,
		image.Point{},
		draw.Src,
	)
	sdi := createSubdividedImage(completeMockImg , tileSize)
	b := color.RGBA{0, 0, 255, 255}
	if sdi[0][0] != b {
		t.Fatalf("not blue")
	}
	r := color.RGBA{255, 0, 0, 255}
	if sdi[0][1] != r {
		t.Fatalf("not red")
	}
	g := color.RGBA{0, 255, 0, 255}
	if sdi[1][0] != g {
		t.Fatalf("not green")
	}
	blk := color.RGBA{0, 0, 0, 255}
	if sdi[1][1] != blk {
		t.Fatalf("not black")
	}
}

func TestBadFoldeForCreateTileImages(t *testing.T) {
	err := os.Mkdir("temp", 0777)
	if err != nil {
		t.Fatalf("could not create test folder")
	}
	os.Create("temp/test.txt")
	os.Create("temp/test2.txt")
	wantLength := 0
	gotLength := len(createTileImages("temp", 5565))
	if wantLength != gotLength {
		t.Fatalf("want %d got %d", wantLength, gotLength)
	}
	err = os.RemoveAll("temp")
	if err != nil {
		t.Fatalf("could not rm dir")
	}
}

func TestNearestRGB(t *testing.T) {
	rgb := color.RGBA{255, 255, 255, 255}
	closer := color.RGBA{230, 230, 230, 255}
	farther := color.RGBA{0, 0, 0, 255}
	got := nearestRGB(rgb, map[color.RGBA]image.Image{closer: nil, farther: nil})
	if closer != got {
		t.Fatalf("want %v got %v", closer, got)
	}
} 

func TestDistance(t *testing.T) {
	rgb := color.RGBA {255, 255, 255, 255}
	closer := color.RGBA {254, 253, 189, 255}
	farther := color.RGBA {1, 2, 3, 255}

	if rgbDistance(rgb, closer) > rgbDistance(rgb, farther) {
		t.Fatalf("here")
	}
}
