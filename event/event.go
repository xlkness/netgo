package event

type EventType int32

var (
	EventTypeSocketConnect EventType = 1
	EventTypeSocketClose   EventType = 2
	EventTypeSocketRecv    EventType = 3
)

type Event struct {
	Type EventType
	Msg  struct {
		Tag     uint32
		Payload []byte
	}
}

func NewEvent(t EventType, tag uint32, payload []byte) *Event {
	e := &Event{
		Type: t,
	}
	e.Msg.Tag = tag
	e.Msg.Payload = payload
	return e
}
