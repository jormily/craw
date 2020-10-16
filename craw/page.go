package craw

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"errors"
	"net/url"
	"strconv"
	"strings"

	"../conf"
	_net "../net"
)

//type CPageItem struct {
//	title		string
//	url		string
//	time		string
//}

type CHref struct {
	title		string
	url			string
	time		string
}

type CPage struct {
	url 		string
	eventvalId 	string
	viewState	string
	val 		map[string]string
	state 		int8
	values      url.Values
	postReq		*_net.CHttpRequest
	crawArray	[]interface{}
}

func NewCPage(pageUrl string) *CPage {
	this := new(CPage)
	this.url = pageUrl
	this.val = make(map[string]string)
	this.state = 0
	this.eventvalId = InitialEventvalId
	this.viewState = InitialViewState
	this.values = url.Values{}

	this.crawArray = []interface{}{}
	this.postReq = _net.NewCHttpRequest(this.url)
	this.postReq.SetHeader("Accept", "*/*")
	this.postReq.SetHeader("Accept-Encoding", "gzip, deflate, br")
	this.postReq.SetHeader("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	this.postReq.SetHeader("Cache-Control", "no-cache")
	this.postReq.SetHeader("Connection", "keep-alive")
	this.postReq.SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// 这个后面三个需要
	this.postReq.SetHeader("X-Requested-With", "XMLHttpRequest")
	this.postReq.SetHeader("X-MicrosoftAjax", "DeitemArraya=true")
	this.postReq.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:77.0) Gecko/20100101 Firefox/77.0")

	return this
}


func (this *CPage) ResetState(){
	this.eventvalId = InitialEventvalId
	this.viewState = InitialViewState
}

func (this *CPage) UpdateValue() {
	for key,val := range this.val {
		this.values.Set(key, val)
	}
	this.values.Set("__VIEWSTATE", this.viewState)
	this.values.Set("__EVENTVALIDATION", this.eventvalId)
}

func (this *CPage) SetVal(key string,val string) {
	this.val[key] = val
	this.UpdateValue()
}

func (this *CPage) SetValues(values url.Values) {
	this.values = values
	this.UpdateValue()
}

func (this *CPage) Update(str string) {
	lines := strings.Split(str, "\n")
	lines = strings.Split(lines[len(lines)-1], "|")

	for i := 1; i < len(lines)/4; i++ {
		if lines[(i-1)*4+3] == "__VIEWSTATE" {
			this.viewState = lines[(i-1)*4+4]
		} else if lines[(i-1)*4+3] == "__EVENTVALIDATION" {
			this.eventvalId = lines[(i-1)*4+4]
		}
	}

	this.UpdateValue()
}

func (this *CPage) CrawEx() ([]*CHref,bool) {
	if this.state == 0 {
		this.SetValues(conf.GetConfigValues(1))
		this.state = this.state + 1
	}else if this.state == 1{
		this.SetValues(conf.GetConfigValues(2))
		this.state = this.state + 1
	}

	doc, body, err := _net.HttpPostEx(this.postReq, this.url, this.values.Encode())
	if err != nil {
		return []*CHref{},false
	}

	hrefArray := []*CHref{}
	doc.Find("div#contentlist div.row").Each(func(i int, s *goquery.Selection) {
		href := new(CHref)
		href.url, _ = s.Find("#linkbtnSrc").Eq(0).Attr("href")
		href.title = s.Find("#linkbtnSrc").Eq(0).Text()
		href.time = s.Find(".publishtime").Eq(0).Text()
		hrefArray = append(hrefArray,href)
	})

	this.Update(body)

	pageStr := doc.Find(".showpagecontent #LabelPage").Eq(0).Text()
	pageArray := strings.Split(pageStr, "/")
	currPage, _ := strconv.ParseInt(pageArray[0], 0, 32)
	maxPage, _ := strconv.ParseInt(pageArray[1], 0, 32)

	fmt.Printf("%d--%d \n",currPage,maxPage)
	if currPage == maxPage {
		return hrefArray,true
	}

	return hrefArray,false
}




