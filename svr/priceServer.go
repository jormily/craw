package svr

import (
	"../craw"
	"../erl"
	"fmt"
)


type PriceServer struct {
	erl.BaseServer
	*craw.CPriceCraw
}

func NewPriceServer() *PriceServer{
	this := new(PriceServer)
	this.BaseServer = erl.BaseServer{}
	this.CPriceCraw = craw.NewCPriceCraw()
	return this
}

func (this *PriceServer) HandleCast(req erl.Req){
	switch req.Method {
	case "craw":
		if ar,ok := req.MsgData.([]*craw.CHref);ok {
			for _,h := range ar {
				this.Craw(h)
			}
		}
	default:
		fmt.Errorf("NoticeServer have not cast of %s", req.Method)
	}
}