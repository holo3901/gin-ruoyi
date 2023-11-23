package erweima

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/png"
	"net/http"
)

func Erweima(wangzhi string, i int) ([]byte, error) {
	// 生成二维码
	qrCode, err := qrcode.Encode(wangzhi, qrcode.Highest, 512)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return nil, err
	}
	a := make(map[int]string, 0)
	a[0] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172335.png"
	a[1] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172430.png"
	a[2] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172449.png"
	a[3] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172513.png"
	a[5] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172510.png"
	a[6] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172516.png"
	a[7] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172519.png"
	a[8] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172522.png"
	a[9] = "http://holoforever.fun/QQ%E5%9B%BE%E7%89%8720231115172526.png"
	// 读取背景图
	backgroundImgURL := a[i]
	response, err := http.Get(backgroundImgURL)
	if err != nil {
		fmt.Println("Error fetching background image:", err)
		return nil, err
	}
	defer response.Body.Close()

	backgroundImg, err := imaging.Decode(response.Body)
	if err != nil {
		fmt.Println("Error decoding background image:", err)
		return nil, err
	}

	// 计算二维码在背景图中的位置
	qrX := (backgroundImg.Bounds().Size().X - 512) / 2 // 二维码宽度为256
	qrY := (backgroundImg.Bounds().Size().Y - 512) / 2 // 二维码高度为256

	// 创建一个新的图像，将背景图绘制到其中
	newImg := image.NewRGBA(image.Rect(0, 0, backgroundImg.Bounds().Dx(), backgroundImg.Bounds().Dy()))
	draw.Draw(newImg, newImg.Bounds(), backgroundImg, image.Point{0, 0}, draw.Src)

	// 将二维码绘制到背景图上
	qrImg, err := png.Decode(bytes.NewReader(qrCode))
	if err != nil {
		fmt.Println("Error decoding QR code image:", err)
		return nil, err
	}

	draw.Draw(newImg, qrImg.Bounds().Add(image.Point{qrX, qrY}), qrImg, image.Point{0, 0}, draw.Over)

	/*//保存最终图像在本地
		outputFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Error creating output image file:", err)
		return nil, err
	}
	defer outputFile.Close()
	err = png.Encode(outputFile, newImg)
	if err != nil {
		fmt.Println("Error encoding output image:", err)
		return nil, err
	}*/

	var buf bytes.Buffer
	err = png.Encode(&buf, newImg)
	if err != nil {
		return nil, fmt.Errorf("Error encoding final image: %v", err)
	}

	// 保存最终图像

	return buf.Bytes(), nil

}
