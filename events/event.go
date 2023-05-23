package events

type Event struct {
	Type   string
	Host   IEventDispatcher
	Target any
	Data   any
}
