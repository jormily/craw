package erl

import (
"time"
)


type Command int8

const (
	call Command = iota
	cast
	timer
	shutdown
)

type Req struct {
	Method  	string
	MsgData 	interface{}
	command 	Command
	time 		int32
	replyCh 	chan *Reply
}

type Reply interface{}

type Server struct {
	id 			int32
	name		string
	callback 	IGenServer
	reqCh   	chan *Req
	stop 		bool
	mgr 		*Manager
}

func NewSvr(name string,mgr *Manager,callBack IGenServer) *Server {
	this := new(Server)
	this.name = name
	this.callback = callBack
	this.reqCh =  make(chan *Req,1024)
	this.stop = false
	this.mgr = mgr

	return this
}

func (this *Server) loop() {
	defer func(){
		this.mgr.Done()
	}()

	for {
		req := <-this.reqCh
		switch req.command {
		case call:
			reply := this.callback.HandleCall(*req)
			req.replyCh<-&reply
		case cast:
			this.callback.HandleCast(*req)
		case timer:
			this.callback.HandleInfo(*req)
			//time.AfterFunc(time.Duration(req.time) * time.Second, func() {
			//  this.Callback.HandleInfo(req)
			//})
		case shutdown:
			this.stop = true
			close(this.reqCh)
			this.callback.OnStop()
			return
		}
	}
}

func (this *Server) Call(method string, msg interface{}) (reply Reply) {
	req := &Req{Method: method, MsgData: msg, command: call, replyCh: make(chan *Reply,1) }
	this.reqCh <-req
	reply = <-req.replyCh
	close(req.replyCh)
	return
}

func (this *Server) Cast(method string, msg interface{}) {
	this.reqCh <- &Req{Method: method, MsgData: msg, command: cast}
}

func (this *Server) SendAfter(method string, seconds int32, msg interface{}) {
	time.AfterFunc(time.Duration(seconds) * time.Second, func() {
		if !this.stop {
			this.reqCh <- &Req{Method: method, MsgData: msg, command: timer, time:seconds}
		}
	})
}

func (this *Server) Stop() {
	this.reqCh <- &Req{command: shutdown}
}
