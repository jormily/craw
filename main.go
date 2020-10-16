package main

import "./analysis"

func main() {
	pa := analysis.NewCPriceAnalysis()
	pa.ImportExcel(".\\土地拍卖价格.xlsx")
	pa.ExportExcel(".\\土地价格.xlsx")
}

