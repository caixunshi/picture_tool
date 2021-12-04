package main

import (
	"picture_tool/utils"
)

/**
* @Author: sirong.huang
* @Date: 2021/12/4 7:52 下午
 */
func main() {
	// 原文件地址
	srcPath := "/Users/shi.cai/Downloads/智能分词前后端交互时序图.png"
	// 输出比例
	scale := 0.7
	utils.GeneratePicture(srcPath, scale)
}
