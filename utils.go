package qrcode

import (
	"image"
	"image/color"
	"math"
)

func NewCircleMask(img image.Image, p image.Point, d int) CircleMask {
	return CircleMask{img, p, d}
}

func DefaultCircleMask(img image.Image) CircleMask {
	w := img.Bounds().Max.X - img.Bounds().Min.X
	h := img.Bounds().Max.Y - img.Bounds().Min.Y

	d := w
	if w > h {
		d = h
	}
	return CircleMask{img, image.Point{0, 0}, d / 2}
}

type CircleMask struct {
	image    image.Image
	point    image.Point
	diameter int
}

func (ci CircleMask) ColorModel() color.Model {
	return ci.image.ColorModel()
}

func (ci CircleMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, ci.diameter, ci.diameter)
}

func (ci CircleMask) At(x, y int) color.Color {
	d := ci.diameter
	dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
	if dis > float64(d)/2 {
		// return ci.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
		return color.Alpha16{0}
	} else {
		return ci.image.At(ci.point.X+x, ci.point.Y+y)
	}
}

func NewRoundMask(img image.Image, r int) RoundMask {
	return RoundMask{
		image: img,
		round: r,
	}
}

type RoundMask struct {
	image image.Image
	round int
}

func (ri RoundMask) ColorModel() color.Model {
	return ri.image.ColorModel()
}

func (ri RoundMask) Bounds() image.Rectangle {
	return ri.image.Bounds()
}

func (ri RoundMask) At(x, y int) color.Color {
	b := ri.image.Bounds()
	w := b.Dx()
	h := b.Dy()
	r := ri.round

	p1 := image.Point{r, r}
	p2 := image.Point{w - r, r}
	p3 := image.Point{r, h - r}
	p4 := image.Point{w - r, h - r}
	if (p1.X-x)*(p1.X-x)+(p1.Y-y)*(p1.Y-y) > r*r && x <= p1.X && y <= p1.Y {
		// return ri.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
		return color.Alpha16{0}
	} else if (p2.X-x)*(p2.X-x)+(p2.Y-y)*(p2.Y-y) > r*r && x > p2.X && y <= p2.Y {
		// return ri.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
		return color.Alpha16{0}
	} else if (p3.X-x)*(p3.X-x)+(p3.Y-y)*(p3.Y-y) > r*r && x <= p3.X && y > p3.Y {
		// return ri.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
		return color.Alpha16{0}
	} else if (p4.X-x)*(p4.X-x)+(p4.Y-y)*(p4.Y-y) > r*r && x > p4.X && y > p4.Y {
		// return ri.image.ColorModel().Convert(color.RGBA{255, 255, 255, 0})
		return color.Alpha16{0}
	} else {
		return ri.image.At(x, y)
	}
}
