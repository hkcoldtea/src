package main

import (
	"gocv.io/x/gocv"
)

func Clarity(mat gocv.Mat) (float64, float64) {
	if mat.Empty() {
		return 0.0, 0.0
	}

	matClone := mat.Clone()
	defer matClone.Close()

	//如果圖片是多通道 就進去轉換
	if mat.Channels() != 1 {
		//將圖像轉換爲灰度顯示
		gocv.CvtColor(mat, &matClone, gocv.ColorRGBToGray)
	}

	destCanny := gocv.NewMat()
	defer destCanny.Close()

	//邊緣檢測
	gocv.Canny(matClone, &destCanny, 200, 200)

	destCannyC := gocv.NewMat()
	defer destCannyC.Close()
	destCannyD := gocv.NewMat()
	defer destCannyD.Close()

	//求矩陣的均值與標準差
	gocv.MeanStdDev(destCanny, &destCannyC, &destCannyD)
	if destCannyD.GetDoubleAt(0, 0) == 0 {
		return 0.0, 0.0
	}

	destA := gocv.NewMat()
	defer destA.Close()

	//Laplace算子
	gocv.Laplacian(matClone, &destA, gocv.MatTypeCV64F, 3, 1, 0, gocv.BorderDefault)

	destC := gocv.NewMat()
	defer destC.Close()
	destD := gocv.NewMat()
	defer destD.Close()

	gocv.MeanStdDev(destA, &destC, &destD)

	destMean := gocv.NewMat()
	defer destMean.Close()

	//Laplace算子
	gocv.Laplacian(matClone, &destMean, gocv.MatTypeCV16U, 3, 1, 0, gocv.BorderDefault)
	mean := destMean.Mean()

	return mean.Val1, destD.GetDoubleAt(0, 0)
}
