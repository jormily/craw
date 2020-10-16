package craw

import (
	"../net"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"reflect"
	"strconv"
	"strings"
)

type CNotic struct {
	CCrawOpt
	Id         	string 	//宗地编号
	Nature   	string 	//宗地性质
	NatureId    int 	//净用地面积
}

type CNoticeCraw struct {
	*CCraw
	natureMap	map[string]int
	crawArray	[]interface{}
}

func NewCNoticeCraw() *CNoticeCraw {
	this := new(CNoticeCraw)
	this.CCraw = NewCCraw(2,nil,)
	this.natureMap = map[string]int{}
	return this
}

func (this *CNoticeCraw) GetNatureMap() map[string]int {
	return this.natureMap
}


func (this *CNoticeCraw) Craw(href *CHref) {
	urls := strings.Split(href.url,"?id=")
	fmt.Printf("craw id = %s \n", urls[len(urls)-1])

	url := "https://www.cdggzy.com/site/LandTrade/ContentHTML.aspx?id=" + urls[len(urls)-1]
	doc,err := net.HttpGet(url)
	if err != nil {
		fmt.Errorf("%s/id = %s",err.Error(),urls[len(urls)-1])
	}

	headMap := map[int]string{}
	doc.Find("table tbody tr").Eq(0).Find("td").Each(
		func(i int, s *goquery.Selection) {
			headMap[i] = s.Text()
		})


	nameVarMap := this.GetNameVarMap()
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}

		if _,err := strconv.ParseInt(s.Find("td").Eq(0).Text(),0,32);err == nil {

			notice := CNotic{}
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				if id, ok := headMap[i]; ok {
					if vid, ok := nameVarMap[id]; ok {
						val := reflect.ValueOf(&notice).Elem()
						val.FieldByName(vid).Set(reflect.ValueOf(s.Text()))
					}
				}
			})

			if strings.Contains(notice.Nature, "住宅") {
				notice.NatureId = 1
			} else {
				notice.NatureId = 0
			}

			this.natureMap[notice.Id] = notice.NatureId
		}
	})
}

