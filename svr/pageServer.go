package svr

import (
	"../erl"
	"../craw"
	"fmt"
)

type IPager interface {
	OnStart()
	CheckStop() bool
}

type PageServer struct {
	erl.BaseServer
	*craw.CPage
	consumer 	*erl.BaseServer
	pager 		IPager
}

func NewPageServer(url string) *PageServer {
	this := new(PageServer)
	this.BaseServer = erl.BaseServer{}
	this.CPage = craw.NewCPage(url)
	return this
}

func (this *PageServer) SetConsumer(c *erl.BaseServer){
	this.consumer = c
}

func (this *PageServer) SetPager(pager IPager){
	this.pager = pager
}

func (this *PageServer) OnStart(){
	this.Cast("craw",nil)
	if this.pager != nil {
		this.pager.OnStart()
	}
}

func (this *PageServer) CheckOver() bool {
	return true
}

func (this *PageServer) HandleCast(req erl.Req){
	switch req.Method {
	case "craw":
		ar,over := this.CrawEx()
		this.consumer.Cast("craw", ar)
		if over && (this.pager == nil || this.pager.CheckStop()) {
			this.Stop()
			this.consumer.Stop()
		}else {
			this.Cast("craw", nil)
		}
	default:
		fmt.Errorf("PageServer have not cast of %s",req.Method)
	}
}

type NoticePager struct {
	*PageServer
	state 			int8
}

func NewNoticePager(p *PageServer) *NoticePager {
	this := new(NoticePager)
	this.PageServer = p
	this.state = 0
	return this
}

func (this *NoticePager) OnStart(){
	this.state = 1
	this.SetVal("ctl00$ContentPlaceHolder1$displaytypeval","7")
}

func (this *NoticePager) CheckStop() bool {
	if this.state == 2 {
		return true
	}

	this.ResetState()
	this.SetVal("ctl00$ContentPlaceHolder1$displaytypeval","8")
	this.state = 2
	return false
}

type PricePager struct {
	*PageServer
}

func NewPricePager(p *PageServer) *PricePager {
	this := new(PricePager)
	this.PageServer = p
	return this
}

func (this *PricePager) OnStart(){
	this.SetVal("ctl00$ContentPlaceHolder1$displaytypeval","2")
}

func (this *PricePager) CheckStop() bool {
	return true
}

