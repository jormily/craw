package net

import (
	"errors"
	"fmt"
	"strings"
	"net/http"
	"io"
)

type CHttpRequest struct {
	url 	string
	header 	map[string]string	
	payload string
}


func NewCHttpRequest(parms ...string) *CHttpRequest {
	r := new(CHttpRequest)
	switch len(parms) {
		case 0:
			r.header = make(map[string]string)
			return r
		case 1:
			r.url = parms[0]
			r.header = make(map[string]string)
			return r	
		case 2:
			r.url = parms[0]	
			r.payload = parms[1]
			r.header = make(map[string]string)
			return r	
	}

	return nil
}

func (r *CHttpRequest) SetUrl(url string) {
	r.url = url
}

func (r *CHttpRequest) SetHeader(key string,val string) {
	r.header[key] = val 
}

func (r *CHttpRequest) SetPayload(payload string) {
	r.payload = payload 
}

func (r *CHttpRequest) Get() (reader io.ReadCloser,err error) {
	res, err := http.Get(r.url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	return res.Body,nil
}

func (r *CHttpRequest) Post() (reader io.ReadCloser,err error) {
	req, err := http.NewRequest("POST", r.url, strings.NewReader(r.payload))
	if err != nil {
		return nil, err
	}

	for key,val := range r.header {
		req.Header.Set(key, val)
	}

	var client = http.Client{}  	//创建客户端
	res, err := client.Do(req) 		//发送请求
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

