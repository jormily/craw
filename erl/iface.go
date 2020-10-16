package erl

type IGenServer interface{
	SetSvr(*Server)
	OnStart()
	HandleCall(Req) Reply
	HandleCast(Req)
	HandleInfo(Req)
	OnStop()
}


