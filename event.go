package server

import "github.com/Sirupsen/logrus"

// EventBus 事件总线
type EventBus struct {
	useQueue bool // 是否使用队列来处理事件
	handlers map[string][]EventHandler
}

// NewEventBus 新建事件总线
func NewEventBus(useQueue bool) *EventBus {
	return &EventBus{
		useQueue: useQueue,
		handlers: make(map[string][]EventHandler),
	}
}

// AddHandler 添加事件处理程序
func (eb *EventBus) AddHandler(eventName string, handler EventHandler) {
	if _, ok := eb.handlers[eventName]; !ok {
		eb.handlers[eventName] = []EventHandler{}
	}
	eb.handlers[eventName] = append(eb.handlers[eventName], handler)
}

// ApplyEvent 触发某个事件
func (eb *EventBus) ApplyEvent(context *Context, e Event) {
	handlers, ok := eb.handlers[e.Name()]
	if !ok || len(handlers) == 0 {
		return
	}
	go func() {
		for _, handler := range handlers {
			err := handler(context, e)
			if err != nil {
				logrus.Error("处理事件:", e.Name(), err)
			}
		}
	}()
}

// Event 事件接口
type Event interface {
	Serialize() []byte
	Name() string
}

// EventHandler 事件处理器
type EventHandler func(*Context, Event) error
