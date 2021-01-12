package qrcode

import (
	"image"
	"image/color"
	"io"

	"github.com/skip2/go-qrcode"
)

const (
	// Level L: 7% error recovery.
	Low = qrcode.Low

	// Level M: 15% error recovery. Good default choice.
	Medium = qrcode.Medium

	// Level Q: 25% error recovery.
	High = qrcode.High

	// Level H: 30% error recovery.
	Highest = qrcode.Highest
)

type IImage interface {
	Create(*qrcode.QRCode, image.Image) (image.Image, error)
}

type IQrCode interface {
	// 设置生成图像圆角
	SetRound(int)

	// 设置头像
	SetAvatar(IImage)

	// 设置背景图
	SetBackgroundImage(IImage)

	// 设置背景颜色
	SetBackgroundColor(color.Color)

	// 设置前景图
	SetForegroundImage(IImage)

	// 设置前景颜色
	SetForegroundColor(color.Color)

	DisableBorder(bool)

	// 返回生成的二维码图片字节数组
	Bytes(size int) ([]byte, error)

	// 返回生成的二维码图片
	Image(size int) (image.Image, error)

	// 将二维码以PNG写入io.Writer
	Write(size int, out io.Writer) error

	// 将二维码以PNG写入指定的文件
	WriteFile(size int, filename string) error
}
