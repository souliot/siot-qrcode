package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestQRCode(t *testing.T) {
	qr, err := New("https://www.baidu.com", Medium)
	if err != nil {
		t.Fatal(err)
	}
	qr.SetRound(20)
	// qr.SetBackgroundColor(&color.RGBA{46, 216, 123, 255})
	qr.SetForegroundColor(&color.RGBA{210, 10, 10, 128})
	// qr.DisableBorder(true)
	qr.SetAvatar(&Avatar{
		Src:    "./data/1.jpg",
		Width:  60,
		Height: 60,
		Round:  8,
	})
	err = qr.WriteFile(256, "./data/out-1.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCircleMask(t *testing.T) {
	file, err := os.Create("./data/1.png")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageFile, err := os.Open("./data/1.jpg")

	if err != nil {
		fmt.Println(err)
	}
	defer imageFile.Close()

	srcImg, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Println(err)
	}

	w := srcImg.Bounds().Max.X - srcImg.Bounds().Min.X
	h := srcImg.Bounds().Max.Y - srcImg.Bounds().Min.Y

	d := w
	if w > h {
		d = h
	}

	dstImg := NewCircleMask(srcImg, image.Point{0, 0}, d)

	png.Encode(file, dstImg)
}

func TestRoundCircleMask(t *testing.T) {
	file, err := os.Create("./data/2.png")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageFile, err := os.Open("./data/1.jpg")

	if err != nil {
		fmt.Println(err)
	}
	defer imageFile.Close()

	srcImg, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Println(err)
	}

	w := srcImg.Bounds().Max.X - srcImg.Bounds().Min.X

	dstImg := NewRoundMask(srcImg, w/2)

	png.Encode(file, dstImg)
}

func TestRoundMask(t *testing.T) {
	file, err := os.Create("./data/3.png")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageFile, err := os.Open("./data/1.jpg")

	if err != nil {
		fmt.Println(err)
	}
	defer imageFile.Close()

	srcImg, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Println(err)
	}

	// w := srcImg.Bounds().Max.X - srcImg.Bounds().Min.X

	dstImg := NewRoundMask(srcImg, 20)

	png.Encode(file, dstImg)
}
