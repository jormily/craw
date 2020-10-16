package anysc

import "sync"

type CManager struct {
	chanelMap 	map[string]chan *CMessage
	wg			sync.WaitGroup
}

func NewCManager() *CManager {
	this := new(CManager)
	this.chanelMap = make(map[string]chan *CMessage)
	this.wg = sync.WaitGroup{}

	return this
}

func (this *CManager)NewConsumer(ident string,consumer IConsumer) *CConsumer {
	if chanel,ok := this.chanelMap[ident];ok {
		c := NewCConsumer(this,consumer,chanel,ident)
		this.AddCunsumer(c)
		return c
	}else {
		chanel := new(chan *CMessage)
		this.chanelMap[ident] = *chanel
		c := NewCConsumer(this,consumer,*chanel,ident)
		this.AddCunsumer(c)
		return c
	}
}

func (this *CManager) NewProducer(ident string,producer IProducer) *CProducer {
	if chanel,ok := this.chanelMap[ident];ok {
		p := NewCProducer(this,producer,chanel,ident)
		this.AddProducer(p)
		return p
	}else{
		chanel := new(chan *CMessage)
		this.chanelMap[ident] = *chanel
		p := NewCProducer(this,producer,*chanel,ident)
		this.AddProducer(p)
		return p
	}
}

func (this *CManager) AddCunsumer(c *CConsumer){
	go this.Cunsumer(c)
	this.wg.Add(1)
}

func (this *CManager) AddProducer(p *CProducer){
	go this.Producer(p)
	this.wg.Add(1)
}

func (this *CManager) Cunsumer(c *CConsumer){
	for {
		select {
		case msg,ok := <-c.messages :
			if !ok || msg.State == MSG_STAST_CLOSE {
				goto Do
			}

			c.consumer.Consume(msg)
		default:


		}
	}

	Do:
		this.wg.Done()
}

func (this *CManager) Producer(p *CProducer){
	for {
		messgs := p.producer.Produce()
		if len(messgs) == 0 || messgs[len(messgs)-1].State == MSG_STAST_CLOSE {
			goto Do
		}

		for _,msg := range messgs {
			p.messages<-msg
		}
	}

	Do:
		this.wg.Done()
}


