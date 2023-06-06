package events

import (
	"fmt"
	"sync"
)

type IEventDispatcher interface {
	DispatchEvent(event *Event)
	AddEventListener(eventType string, listener func(*Event), target interface{})
	RemoveEventListener(eventType string, listener func(*Event), target interface{})
	HasEventListener(eventType string) bool
}

type EventDispatcher struct {
	sync.RWMutex

	host interface{}
	list map[string]*dispatcher
}

func (this *EventDispatcher) Constructor(host interface{}) {
	if this != nil {
		this.host = host
		this.list = make(map[string]*dispatcher)
	}
}
func (this *EventDispatcher) DispatchEvent(event *Event) {
	if this != nil && this.list != nil {
		event.Host = this
		if d, ok := this.list[event.Type]; ok {
			if event.Target == nil {
				if this.host != nil {
					event.Target = this.host
				} else {
					event.Target = this
				}
			}
			d.dispatch(event)
		}
	}
}
func (this *EventDispatcher) AddEventListener(eventType string, listener func(*Event), target interface{}) {
	if this != nil {
		if this.list == nil {
			this.list = make(map[string]*dispatcher)
		}
		if d, ok := this.list[eventType]; ok {
			d.addListener(listener, target)
		} else {
			this.Lock()
			defer this.Unlock()
			dph := new(dispatcher)
			dph.Constructor(eventType)
			dph.addListener(listener, target)
			this.list[eventType] = dph
		}
	}
}
func (this *EventDispatcher) RemoveEventListener(eventType string, listener func(*Event), target interface{}) {
	if this != nil && this.list != nil {
		if d, ok := this.list[eventType]; ok {
			this.Lock()
			defer this.Unlock()
			d.removeListener(listener, target)
		}
	}
}
func (this *EventDispatcher) RemoveAllEventListener() {
	if this != nil && this.list != nil {
		this.Lock()
		defer this.Unlock()
		this.list = make(map[string]*dispatcher)
	}
}
func (this *EventDispatcher) HasEventListener(eventType string) bool {
	if this != nil && this.list != nil {
		if _, ok := this.list[eventType]; ok {
			return true
		}
	}
	return false
}

//++++++++++++++++++++ dispatch ++++++++++++++++++++

type dispatcher struct {
	sync.RWMutex

	Type     string
	listener map[string]func(*Event)
}

func (this *dispatcher) Constructor(eventType string) {
	this.Type = eventType
	this.listener = make(map[string]func(*Event))
}
func (this *dispatcher) dispatch(event *Event) {
	if this.listener == nil {
		return
	}
	for _, lis := range this.listener {
		go lis(event)
	}
}
func (this *dispatcher) addListener(listener func(*Event), target interface{}) {
	if this.listener == nil {
		return
	}
	id := this.createID(listener, target)
	if _, has := this.listener[id]; !has {
		this.Lock()
		defer this.Unlock()
		this.listener[id] = listener
	}
}
func (this *dispatcher) removeListener(listener func(*Event), target interface{}) {
	if this.listener == nil {
		return
	}
	id := this.createID(listener, target)
	if _, has := this.listener[id]; has {
		this.Lock()
		defer this.Unlock()
		delete(this.listener, id)
	}
}
func (this *dispatcher) createID(listener func(*Event), target interface{}) string {
	id := fmt.Sprintf("%d", listener)
	id += "_"
	id += fmt.Sprintf("%d", target)
	return id
}
func (this *dispatcher) dispose() {
	this.listener = nil
}
