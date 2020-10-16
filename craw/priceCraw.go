package craw

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"reflect"
	"strconv"
	"strings"

	"../net"
)

type CPrice struct {
	CCrawOpt
	Id         	string 	//宗地编号
	Location   	string 	//宗地位置
	Measure    	string 	//净用地面积
	InitPrice  	string 	//起始价
	FinalPrice 	string 	//成交价
	Price		string	//最终价格
	Time     	string	//成交时间
	Clincher   	string 	//竞得人
	PageId     	string 	//页面id
	Nature		string  //
}


type CPriceCraw struct {
	*CCraw
	noticeCraw *CNoticeCraw
}

func NewCPriceCraw() *CPriceCraw {
	this := new(CPriceCraw)
	this.CCraw = NewCCraw(1,this)
	return this
}

func (this *CPriceCraw)ExportFactor(item interface{}) bool {
	if price,ok := item.(*CPrice);ok {
		if price.CCrawOpt.exp == 1 {
			return true
		}
	}

	return false
}

func (this *CPriceCraw) Check(nc *CNoticeCraw) {
	natureMap := nc.GetNatureMap()
	for _,pr := range this.itemArray {
		if price,ok := pr.(*CPrice);ok {
			if id, ok := natureMap[price.Id]; ok {
				if id == 1 {
					price.Nature = "住宅"
					price.CCrawOpt.exp = 1
				} else {
					price.Nature = "其他"
					price.CCrawOpt.exp = 2
				}
			}
			if price.Clincher == "流拍" || price.Clincher == "终止" {
				price.CCrawOpt.exp = 3
			}
		}
	}
}


func (this *CPriceCraw) Craw(href *CHref) {
	urls := strings.Split(href.url,"?id=")
	fmt.Printf("craw id = %s \n", urls[len(urls)-1])

	url := "https://www.cdggzy.com/site/LandTrade/ContentHTML.aspx?id=" + urls[len(urls)-1]
	doc,err := net.HttpGet(url)
	if err != nil {
		fmt.Errorf("%s/id = %s",err.Error(),urls[len(urls)-1])
	}

	headMap := map[int]string{}
	doc.Find("form#form1 tbody tr").Eq(0).Find("th").Each(
		func(i int, s *goquery.Selection) {
			headMap[i] = s.Text()
	})

	nameVarMap := this.GetNameVarMap()
	doc.Find("form#form1 tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}

		price := CPrice{}
		price.PageId = urls[len(urls)-1]
		price.Time = href.time
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if id,ok := headMap[i]; ok {
				if vid, ok := nameVarMap[id]; ok {
					val := reflect.ValueOf(&price).Elem()
					val.FieldByName(vid).Set(reflect.ValueOf(s.Text()))
				}
			}
		})

		if !(price.Clincher == "流拍" || price.Clincher == "终止") {
			if strings.HasSuffix(price.FinalPrice, "万元/亩") {
				price.Price = strings.Replace(price.FinalPrice, "万元/亩", "", 1)
				pr, _ := strconv.ParseFloat(price.Price, 32)
				pr = pr * 15
				price.Price = fmt.Sprintf("%v元/平方米", pr)
			} else if strings.HasSuffix(price.FinalPrice, "元/平方米") {
				price.Price = price.FinalPrice
			}
		}

		this.itemArray = append(this.itemArray, &price)
	})
}
