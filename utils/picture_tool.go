package utils

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	formatStr             = "%s/%d_高=%d_宽=%d.png"
	basePath              = "./file"
	beforePath            = "./file/翻页前的图片"
	afterPathWidthChange  = "./file/单边变动（仅宽变动）"
	afterPathHeightChange = "./file/单边变动（仅高变动）"
	afterPathAllChange    = "./file/双边变动"
)

// GeneratePicture 生成图片
func GeneratePicture(filename string, scale float64) {
	// 检查文件夹是否存在，如果不存在则先创建文件夹
	err := checkFolder()
	if err != nil {
		log.Printf("生成文件目录失败，err:%v", err)
		return
	}

	// 加载图片，生成image对象
	src, err := loadImage(filename)
	if err != nil {
		log.Printf("获取图片失败，err:%v", err)
		return
	}
	// 生成宽和高的比例
	temp := (1 - scale) / 2
	max := scale + temp
	min := scale - temp
	mid := scale

	// 获取原始图片的宽和高
	width := float64(src.Bounds().Max.X)
	height := float64(src.Bounds().Max.Y)

	// 输出翻页前的图
	widthArray := []int{int(min * width), int(mid * width), int(max * width)}
	heightArray := []int{int(min * height), int(mid * height), int(max * height)}
	generate(src, beforePath, widthArray, heightArray)
	// 输出翻页后的图
	// 单边变动（仅宽变动）
	widthArray = []int{
		int(min * width),
		int(mid * width),
		int(max * width),
		int(max * width * min),
		int(max * width),
		int(max * width * max),
	}
	heightArray = []int{
		int(mid * height),
	}
	generate(src, afterPathWidthChange, widthArray, heightArray)

	// 输出翻页后的图
	// 单边变动（仅高变动）
	widthArray = []int{
		int(mid * width),
	}
	heightArray = []int{
		int(min * height),
		int(mid * height),
		int(max * height),
		int(max * height * min),
		int(max * height * mid),
		int(max * height * max),
	}
	generate(src, afterPathHeightChange, widthArray, heightArray)

	// 双边变动
	widthArray = []int{
		int(min * width),
		int(mid * width),
		int(max * width),
		int(max * width * min),
		int(max * width),
		int(max * width * max),
	}
	heightArray = []int{
		int(min * height),
		int(mid * height),
		int(max * height),
		int(max * height * min),
		int(max * height * mid),
		int(max * height * max),
	}
	generate(src, afterPathAllChange, widthArray, heightArray)
}

func checkFolder() error {
	_, err := os.Stat(basePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(basePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(beforePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(beforePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(afterPathWidthChange)
	if os.IsNotExist(err) {
		err = os.Mkdir(afterPathWidthChange, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(afterPathHeightChange)
	if os.IsNotExist(err) {
		err = os.Mkdir(afterPathHeightChange, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(afterPathAllChange)
	if os.IsNotExist(err) {
		err = os.Mkdir(afterPathAllChange, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func generate(src image.Image, filepath string, widthArr, heightArr []int) {
	// index用来记录文件序号
	index := 1
	for _, width := range widthArr {
		for _, height := range heightArr {
			filename := fmt.Sprintf(formatStr, filepath, index, width, height)
			trimming(src, filename, width, height)
			index++
		}
	}
}

// trimming 裁剪图片
func trimming(src image.Image, afterFilename string, w, h int) {
	img, err := imageCopy(src, 1, 1, w, h)
	if err != nil {
		log.Println("image copy fail...")
	}
	saveErr := saveImage(afterFilename, img)
	if saveErr != nil {
		log.Println("save image fail..")
	}
}

// loadImage 加载图片
func loadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

// saveImage 保存图片到指定路径
func saveImage(p string, src image.Image) error {
	f, err := os.OpenFile(p, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(p)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		err = jpeg.Encode(f, src, &jpeg.Options{Quality: 80})
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, src, &gif.Options{NumColors: 256})
	}
	return err
}

// imageCopy 复制图片
func imageCopy(src image.Image, x, y, w, h int) (image.Image, error) {
	var subImg image.Image
	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {
		return subImg, errors.New("图片解码失败")
	}
	return subImg, nil
}
