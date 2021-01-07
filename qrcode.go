package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"

	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
)

type Avatar struct {
	Src    string // 头像地址
	Width  int    // 头像宽度
	Height int    // 头像高度
	Round  int    // 头像圆角
}

func (a *Avatar) Create(qr *qrcode.QRCode, img image.Image) (image.Image, error) {
	avatar, err := os.Open(a.Src)
	if err != nil {
		return nil, fmt.Errorf("open avatar file error: %s", err.Error())
	}

	defer avatar.Close()

	decode, _, err := image.Decode(avatar)
	if err != nil {
		return nil, err
	}

	decode = resize.Resize(uint(a.Width), uint(a.Height), decode, resize.Lanczos3)

	decode = NewRoundMask(decode, a.Round)

	b := img.Bounds()

	// 设置为居中
	offset := image.Pt((b.Max.X-decode.Bounds().Max.X)/2, (b.Max.Y-decode.Bounds().Max.Y)/2)

	m := image.NewRGBA(b)

	draw.Draw(m, b, img, image.Point{X: 0, Y: 0}, draw.Src)

	draw.Draw(m, decode.Bounds().Add(offset), decode, image.Point{X: 0, Y: 0}, draw.Over)

	return m, err
}

type BackgroundImage struct {
	Src    string
	X      int
	Y      int
	Width  int
	Height int
}

func (a *BackgroundImage) Create(qr *qrcode.QRCode, img image.Image) (image.Image, error) {
	file, err := os.Open(a.Src)
	if err != nil {
		return nil, fmt.Errorf("打开背景图文件失败 %s", err.Error())
	}

	img = resize.Resize(uint(a.Width), uint(a.Height), img, resize.Lanczos3)

	defer file.Close()

	bg, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	offset := image.Pt(a.X, a.Y)

	b := bg.Bounds()

	m := image.NewRGBA(b)

	draw.Draw(m, b, bg, image.Point{X: 0, Y: 0}, draw.Src)

	draw.Draw(m, img.Bounds().Add(offset), img, image.Point{X: 0, Y: 0}, draw.Over)

	return m, nil
}

type ForegroundImage struct {
	Src string
}

func (a *ForegroundImage) Create(qr *qrcode.QRCode, img image.Image) (image.Image, error) {
	file, err := os.Open(a.Src)
	if err != nil {
		return nil, fmt.Errorf("打开前景图文件失败 %s", err.Error())
	}

	defer file.Close()

	decode, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	// 获取二维码的宽高
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y

	// 获取要填充的图片宽高
	foregroundW, foregroundH := decode.Bounds().Max.X, decode.Bounds().Max.Y

	if width != foregroundW || height != foregroundH {
		// 如果不一致将填充图剪裁
		decode = resize.Resize(uint(width), uint(height), decode, resize.Lanczos3)
	}

	m := image.NewRGBA(img.Bounds())
	d := image.NewRGBA(decode.Bounds())

	draw.Draw(m, m.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(d, d.Bounds(), decode, image.Point{X: 0, Y: 0}, draw.Src)

	for y := 0; y < img.Bounds().Max.X; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {

			// 检测像素是否为白色或者透明色
			if m.At(x, y).(color.RGBA).R == 255 && m.At(x, y).(color.RGBA).G == 255 && m.At(x, y).(color.RGBA).B == 255 && m.At(x, y).(color.RGBA).A == 255 {
				continue
			}

			if m.At(x, y).(color.RGBA) == qr.BackgroundColor {
				continue
			}

			// 填充颜色
			m.Set(x, y, color.RGBA{R: d.At(x, y).(color.RGBA).R, G: d.At(x, y).(color.RGBA).G, B: d.At(x, y).(color.RGBA).B, A: d.At(x, y).(color.RGBA).A})
		}
	}

	return m, nil
}

type QrCode struct {
	qr              *qrcode.QRCode
	round           int
	Avatar          IImage
	ForegroundImage IImage
	BackgroundImage IImage
}

var _ IQrCode = new(QrCode)

func New(content string, level qrcode.RecoveryLevel) (IQrCode, error) {

	qr, err := qrcode.New(content, level)
	if err != nil {
		return nil, err
	}

	qrCode := &QrCode{}
	qrCode.qr = qr

	return qrCode, nil
}

// 设置生成图像圆角
func (q *QrCode) SetRound(r int) {
	q.round = r
}

// 设置头像
func (q *QrCode) SetAvatar(avatar IImage) {
	q.Avatar = avatar
}

// 设置背景图
func (q *QrCode) SetBackgroundImage(img IImage) {
	q.BackgroundImage = img
}

// 设置背景颜色
func (q *QrCode) SetBackgroundColor(color color.Color) {
	q.qr.BackgroundColor = color
}

// 设置前景图
func (q *QrCode) SetForegroundImage(img IImage) {
	q.ForegroundImage = img
}

// 设置前景颜色
func (q *QrCode) SetForegroundColor(color color.Color) {
	q.qr.ForegroundColor = color
}

func (q *QrCode) DisableBorder(disable bool) {
	q.qr.DisableBorder = disable
}

// 返回生成的二维码图片
func (q *QrCode) Image(size int) (image.Image, error) {
	img := q.qr.Image(size)
	var err error

	if q.ForegroundImage != nil {
		if img, err = q.ForegroundImage.Create(q.qr, img); err != nil {
			return nil, err
		}
	}

	if q.Avatar != nil {
		if img, err = q.Avatar.Create(q.qr, img); err != nil {
			return nil, err
		}
	}

	if q.BackgroundImage != nil {
		if img, err = q.BackgroundImage.Create(q.qr, img); err != nil {
			return nil, err
		}
	}

	img = NewRoundMask(img, q.round)

	return img, nil
}

// 将二维码以PNG写入io.Writer
func (q *QrCode) Write(size int, out io.Writer) error {
	var p []byte

	p, err := q.PNG(size)

	if err != nil {
		return err
	}
	_, err = out.Write(p)
	return err
}

// 将二维码以PNG写入指定的文件
func (q *QrCode) WriteFile(size int, filename string) error {
	var p []byte

	p, err := q.PNG(size)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, p, os.FileMode(0644))
}

// 返回 png 二维码图片
func (q *QrCode) PNG(size int) ([]byte, error) {
	img, err := q.Image(size)
	if err != nil {
		return nil, err
	}
	encoder := png.Encoder{CompressionLevel: png.BestCompression}

	var b bytes.Buffer
	err = encoder.Encode(&b, img)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
