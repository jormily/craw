package analysis

import (
	"../conf"
	"../excel"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type CPrice struct {
	Location 	string
	Price 		string
	Time  		string
}

type CPriceCol struct {
	Time 			string
	TotalPrice  	float64
	Price   		float64
	PriceArray		[]float64
}

type CPriceAnalysis struct {
	excelIn			*excel.CExcel
	excelOut		*excel.CExcel
	sheetMap		map[string][]interface{}
	regionArray 	[]string
}


func (this *CPriceAnalysis)GetSheetMap() map[string][]interface{} {
	return this.sheetMap
}


func NewCPriceAnalysis() *CPriceAnalysis {
	this := new(CPriceAnalysis)
	this.excelIn = excel.NewCExcel(3,nil,reflect.TypeOf(&CPrice{}))
	this.excelOut = excel.NewCExcel(4,nil)
	this.sheetMap = make(map[string][]interface{})
	this.regionArray = []string{}

	for _,r := range conf.GetConfigArray("region") {
		if region, ok := r.(string);ok {
			this.regionArray = append(this.regionArray, region)
		}
	}
	return this
}

func (this *CPriceAnalysis) Add(region string,pr *CPrice){
	ar, ok := this.sheetMap[region];
	if !ok {
		ar = []interface{}{}
		this.sheetMap[region] = ar
	}

	pr.Price = strings.Replace(pr.Price, "元/平方米", "", 1)
	price, _ := strconv.ParseFloat(pr.Price, 16)

	reg := regexp.MustCompile(`^[\d]{4}-[\d]{2}`)

	var pc *CPriceCol = nil
	for _, it := range ar {
		timeStr := reg.FindString(it.(*CPriceCol).Time)
		if timeStr == pr.Time {
			pc = it.(*CPriceCol)
			break
		}
	}

	if pc == nil {
		pc := &CPriceCol{}
		pc.Time = reg.FindString(pr.Time)
		pc.PriceArray = []float64{price}
		pc.TotalPrice = price
		pc.Price = pc.TotalPrice/float64(len(pc.PriceArray))
		ar = append(ar, pc)
		this.sheetMap[region] = ar
	}else {
		pc.PriceArray = append(pc.PriceArray, price)
		pc.TotalPrice = pc.TotalPrice + price
		pc.Price = pc.TotalPrice/float64(len(pc.PriceArray))
	}
}

func (this *CPriceAnalysis) ImportExcel(fileName string) error {
	err := this.excelIn.Import(fileName)
	if err != nil {
		return err
	}

	// 遍历sheet页读取
	fmt.Println(this.excelIn.GetSheetMap())
	//sheetMap := this.excelIn.GetSheeptMap()
	for _, sheet := range *this.excelIn.GetSheetMap() {
		for _, item := range sheet {
			price := item.(*CPrice)
			for _, region := range this.regionArray {
				if strings.Contains(price.Location, region) {
					this.Add(region, price)
				}
			}

		}
	}
	return nil

}

func (this *CPriceAnalysis) ExportExcel(fileName string) {
	this.excelOut.SetSheetMap(this.sheetMap)
	this.excelOut.Export(fileName)
}

