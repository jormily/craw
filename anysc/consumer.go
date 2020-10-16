package anysc

type IConsumer interface {
	Consume(msg *CMessage)
}

type CConsumer struct {
	messages 		<-chan *CMessage
	ident 			string
	consumer		IConsumer
	manager   		*CManager
}


func NewCConsumer(manager *CManager,consumer IConsumer,messages chan *CMessage,ident string) *CConsumer {
	this := new(CConsumer)
	this.manager = manager
	this.consumer = consumer
	this.messages = messages
	this.ident = ident
	return this
}

func (this *CConsumer) SetMessages(messages chan *CMessage){
	this.messages = messages
}

func (this *CConsumer) SetIdent(ident string){
	this.ident = ident
}

func (this *CConsumer) SetConsumer(consumer IConsumer){
	this.consumer = consumer
}

func (this *CConsumer) SetManager(manager *CManager){
	this.manager = manager
}

