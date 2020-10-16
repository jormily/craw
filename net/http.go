package net

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

var (
	PostReq = NewCHttpRequest()
	//InitialViewState = ""
	//InitialEventvalId = ""
)

func HttpGet(url string)  (*goquery.Document, error) {
	r := NewCHttpRequest(url)
	body, err := r.Get()
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func HttpPost(r *CHttpRequest,url string,playload string) (*goquery.Document, error) {
	if r == nil {
		r = NewCHttpRequest(url, playload)
	}else{
		r.SetUrl(url)
		r.SetPayload(playload)
	}
	//PostReq.SetUrl(url)
	//PostReq.SetPayload(playload)
	body, err := r.Post()
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	return doc,nil
}

func HttpPostEx(r *CHttpRequest, url string,playload string) (*goquery.Document, string, error) {
	if r == nil {
		r = NewCHttpRequest(url, playload)
	}else{
		r.SetUrl(url)
		r.SetPayload(playload)
	}

	body, err := r.Post()
	if err != nil {
		return nil, "", err
	}
	defer body.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(body)
	doc, err1 := goquery.NewDocumentFromReader(bytes.NewReader(buffer.Bytes()))
	if err1 != nil {
		return nil, "", err1
	}

	return doc, buffer.String(), nil
}


func init() {
	PostReq.SetHeader("Accept", "*/*")
	PostReq.SetHeader("Accept-Encoding", "gzip, deflate, br")
	PostReq.SetHeader("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	PostReq.SetHeader("Cache-Control", "no-cache")
	PostReq.SetHeader("Connection", "keep-alive")
	PostReq.SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// 这个后面三个需要
	PostReq.SetHeader("X-Requested-With", "XMLHttpRequest")
	PostReq.SetHeader("X-MicrosoftAjax", "DeitemArraya=true")
	PostReq.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:77.0) Gecko/20100101 Firefox/77.0")
}
