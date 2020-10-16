package anysc

var (
	mgr = NewCManager()
)

func NewConsumer(ident string,consumer IConsumer) *CConsumer {
	return mgr.NewConsumer(ident,consumer)
}

func NewProducer(ident string,producer IProducer) *CProducer {
	return mgr.NewProducer(ident,producer)
}

func  AddProducer(p *CProducer){
	mgr.AddProducer(p)
}

func AddConsumer(c *CConsumer){
	mgr.AddCunsumer(c)
}

func Wait(){
	mgr.wg.Wait()
}