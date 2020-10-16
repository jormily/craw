package excel

import (
	"../conf"

	"fmt"
	"github.com/tealeg/xlsx"
	"reflect"
)

type IExcelFactor interface {
	ExportFactor(item interface{}) bool
}

type CExcelVar struct {
	Variable	string
	Name 		string
}

type CExcel struct {
	eType		int8
	sheetMap	map[string][]interface{}
	variArray	[]CExcelVar
	factor 		IExcelFactor
	typ         reflect.Type
}

func NewCExcel(eType int8,args ...interface{}) *CExcel {
	this := new(CExcel)
	this.eType = eType
	this.sheetMap = map[string][]interface{}{}
	this.variArray = []CExcelVar{}
	this.factor = nil
	this.typ = nil

	for i,arg := range args {
		if i == 0 {
			if ef,ok := arg.(IExcelFactor); ok {
				this.factor = ef
			}
		}

		if i == 1 {
			if typ,ok := arg.(reflect.Type); ok {
				this.typ = typ
			}
		}
	}


	if this.init() {
		return this
	}else{
		return nil
	}
}

func (this *CExcel)init() bool {
	config := conf.GetConfig("excel",int(this.eType))
	if config == nil {
		return false
	}

	if cv,ok := config.(map[string]interface{}); ok {
		for k,v := range cv {
			this.variArray = append(this.variArray ,CExcelVar{k,v.(string)})
		}
	}

	return true
}

func (this *CExcel)Add(sheetId string,it interface{}){
	if this.typ != reflect.TypeOf(it) {
		return
	}

	if sheet,ok := this.sheetMap[sheetId];ok {
		this.sheetMap[sheetId] = append(sheet, it)
	}else {
		this.sheetMap[sheetId] = []interface{}{it}
	}
}

func (this *CExcel)SetSheetMap(sheetMap map[string][]interface{}) {
	this.sheetMap = sheetMap
}

func (this *CExcel)SetSheetMapEx(key string,arr []interface{}){
	this.sheetMap[key] = arr
}

func (this *CExcel)GetVarNameMap() map[string]string {
	r := map[string]string{}
	for _,v := range this.variArray {
		r[v.Variable] = v.Name
	}
	return r
}

func (this *CExcel)GetNameVarMap() map[string]string {
	r := map[string]string{}
	for _,v := range this.variArray {
		r[v.Name] = v.Variable
	}
	return r
}

func (this *CExcel)GetSheetMap() *map[string][]interface{} {
	return &this.sheetMap
}


func (this *CExcel)Import(fileName string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		//fmt.Errorf(err.Error())
		return err
	}

	typ := this.typ
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	nameVarMap := this.GetNameVarMap()
	for _, sheet := range file.Sheets {
		varidMap := map[int]string{}
		if row, err := sheet.Row(0); err == nil {
			for i := 0; i < row.Sheet.MaxCol; i++ {
				if varid, ok := nameVarMap[row.GetCell(i).String()]; ok {
					varidMap[i] = varid
				}
			}
		}

		for i := 1; i < sheet.MaxRow; i++ {
			val := reflect.New(typ).Elem()
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			if row, err := sheet.Row(i);err == nil {
				for j := 0; j < sheet.MaxCol; j++ {
					if varid, ok := varidMap[j]; ok {
						val.FieldByName(varid).Set(reflect.ValueOf(row.GetCell(j).String()))
					}
				}

				it := val.Addr().Interface()
				if data, ok := this.sheetMap[sheet.Name]; ok {
					this.sheetMap[sheet.Name] = append(data, it)
				} else {
					data := []interface{}{}
					data = append(data, it)
					this.sheetMap[sheet.Name] = data
				}
			}
		}
	}
	return nil
}


func (this *CExcel)Export(fileName string) error {
	file := xlsx.NewFile()
	for sheetName,sheetData := range this.sheetMap {
		sheet, err := file.AddSheet(sheetName)
		if err != nil {
			fmt.Errorf(err.Error())
			return err
		}

		row := sheet.AddRow()
		row.SetHeightCM(0.5)
		for _, v := range this.variArray {
			cell := row.AddCell()
			cell.Value = v.Name
		}


		for _, item := range sheetData {
			if this.factor == nil || this.factor.ExportFactor(item) {
				row = sheet.AddRow()
				row.SetHeightCM(0.5)

				val := reflect.ValueOf(item).Elem()
				//if val.Kind() == reflect.Ptr {
				//	val = val.Elem()
				//}
				for _, v := range this.variArray {
					cell := row.AddCell()
					cell.Value = fmt.Sprintf("%v", val.FieldByName(v.Variable))
				}
			}
		}
	}

	err := file.Save(fileName)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	return nil
}
