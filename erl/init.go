package erl

var (
	mgr = NewManager()
)

func Register(name string,callBack IGenServer) bool {
	return mgr.Register(name,callBack)
}

func UnRegister(name string) {
	mgr.UnRegister(name)
}

func Get(name string) *Server {
	return mgr.Get(name)
}

func Count() int {
	return mgr.Count()
}

func Stop(name string) (err error) {
	return mgr.Stop(name)
}

func Call(name string,method string, msg interface{}) (reply Reply,err error) {
	return mgr.Call(name,method,msg)
}

func Cast(name string,method string, msg interface{}) (err error) {
	return mgr.Cast(name,method,msg)
}

func SendAfter(name string,method string, seconds int32, msg interface{}) (err error) {
	return mgr.SendAfter(name,method,seconds,msg)
}

func Wait(){
	mgr.Wait()
}
