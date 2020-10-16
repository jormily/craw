package erl

import (
	"fmt"
	"sync"
)


type Manager struct {
	pidMap 	map[string]*Server
	sync.RWMutex
	sync.WaitGroup
}

func NewManager() *Manager {
	this := new(Manager)
	this.pidMap = make(map[string]*Server)
	return this
}

func (this *Manager) Register(name string,callBack IGenServer) bool {
	this.Lock()
	defer this.Unlock()

	if _,ok := this.pidMap[name]; ok {
		return false
	}

	p := NewSvr(name,this,callBack)
	p.callback.SetSvr(p)
	p.callback.OnStart()
	this.Add(1)

	go p.loop()

	return true
}

func (this *Manager) UnRegister(name string) {
	this.Lock()
	defer this.Unlock()
	delete(this.pidMap, name)
}

func (this *Manager) Get(name string) *Server {
	this.Lock()
	defer this.Unlock()
	if p,ok := this.pidMap[name];ok {
		return p
	}
	return nil
}

func (this *Manager) Count() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.pidMap)
}

func (this *Manager) Stop(name string) (err error){
	defer func(){
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	p := this.Get(name)
	if p != nil {
		p.Stop()
	}

	return
}

func (this *Manager) Call(name string,method string, msg interface{}) (reply Reply,err error) {
	defer func(){
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	p := this.Get(name)
	if p != nil {
		reply = p.Call(method,msg)
	}

	return
}


func (this *Manager) Cast(name string,method string, msg interface{}) (err error) {
	defer func(){
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	p := this.Get(name)
	if p != nil {
		p.Cast(method,msg)
	}

	return
}

func (this *Manager) SendAfter(name string,method string, seconds int32, msg interface{}) (err error) {
	defer func(){
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	p := this.Get(name)
	if p != nil {
		p.SendAfter(method,seconds,msg)
	}

	return
}