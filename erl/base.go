package erl

type BaseServer struct {
	*Server
}

func (this *BaseServer) SetSvr(srv *Server) {
	this.Server = srv
}

func (this *BaseServer) OnStart(){}

func (this *BaseServer) OnStop(){}

func (this *BaseServer) HandleCall(req Req) (reply Reply) {return}

func (this *BaseServer) HandleCaset(req Req){}

func (this *BaseServer) HandleInfo(req Req){}
