package svr

import (
	"../craw"
	"../erl"
	"fmt"

	//"fmt"
)


type NoticeServer struct {
	erl.BaseServer
	*craw.CNoticeCraw
}

func NewNoticeServer() *NoticeServer{
	this := new(NoticeServer)
	this.BaseServer = erl.BaseServer{}
	this.CNoticeCraw = craw.NewCNoticeCraw()
	return this
}

func (this *NoticeServer) HandleCast(req erl.Req){
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