package anysc

type IProducer interface {
	Produce() []*CMessage
}

type CProducer struct {
	messages 	chan<- *CMessage
	ident 		string
	manager   	*CManager
	producer 	IProducer
}

func NewCProducer(manager *CManager,producer IProducer,messages chan *CMessage,ident string) *CProducer{
	this := new(CProducer)
	this.ident = ident
	this.manager = manager
	this.producer = producer
	this.messages = messages
	return this
}
