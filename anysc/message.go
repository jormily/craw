package anysc

type MSG_STATE int8

const (
	MSG_STAST_CLOSE MSG_STATE = iota
	MSG_STATE_WORK
)

type CMessage struct {
	State 	MSG_STATE
	Msg 	interface{}
}
