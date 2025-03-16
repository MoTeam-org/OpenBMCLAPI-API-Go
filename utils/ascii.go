package utils

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
)

// 将图片URL转换为ASCII艺术
func ImageToAscii(url string, width int) (string, error) {
	// 下载图片
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取图片数据
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", err
	}

	// 获取图片尺寸
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	// 计算高度，保持宽高比
	height := width * imgHeight / imgWidth

	// 优化的ASCII字符集，从暗到亮
	ascii := []string{" ", ".", ":", "-", "=", "+", "*", "#", "%", "@"}

	var result string
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 映射坐标
			srcX := x * imgWidth / width
			srcY := y * imgHeight / height

			// 获取像素并计算亮度
			pixel := img.At(srcX, srcY)
			gray := color.GrayModel.Convert(pixel).(color.Gray)
			brightness := gray.Y

			// 映射亮度到ASCII字符
			index := int(brightness) * (len(ascii) - 1) / 255
			result += ascii[index]
		}
		result += "\n"
	}

	return result, nil
}
