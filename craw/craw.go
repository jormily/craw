package craw

import (
	//"../conf"
	"../excel"
	//"fmt"
	//"github.com/tealeg/xlsx"
	//"reflect"
)

type ICraw interface {
	Run()

}

type VarName struct {
	Variable 	string
	Name 		string
}

type CCrawOpt struct {
	exp 		int8
}

type CCraw struct {
	*excel.CExcel
	itemArray	[]interface{}
	crawArray	[]interface{}
}

func NewCCraw(eType int8,args ...interface{}) *CCraw {
	this := new(CCraw)
	this.CExcel = excel.NewCExcel(eType, args...)
	this.itemArray = []interface{}{}
	this.crawArray = []interface{}{}
	return this
}

func (this *CCraw)SetCrawArray(arr []interface{}){
	this.crawArray = arr
}

func (this *CCraw)GetCrawArray() []interface{}{
	return this.crawArray
}

func (this *CCraw)GetItemArray() []interface{} {
	return this.itemArray
}

func (this *CCraw)ExportExcel(fileName string) error {
	this.CExcel.SetSheetMapEx("sheet1",this.itemArray)
	return this.CExcel.Export(fileName)
}