package helper

import (
	"sync"
	"time"
)

var (
	notification_ports = map[string]chan interface{}{}
	notification_mutex sync.Mutex
)

func CreateNotificationPort(id string, buffer int) chan interface{} {
	notification_mutex.Lock()
	defer notification_mutex.Unlock()
	port, ok := notification_ports[id]
	if !ok {
		port = make(chan interface{}, buffer)
		notification_ports[id] = port
	}
	return port
}

func GetNotificationPort(id string) chan interface{} {
	notification_mutex.Lock()
	defer notification_mutex.Unlock()
	port, ok := notification_ports[id]
	if !ok {
		port = make(chan interface{}, 1)
		notification_ports[id] = port
	}
	return port
}

func Notify(id string, data interface{}) {
	GetNotificationPort(id) <- data
}

func NotifyTimeout(id string, data interface{}, timeout time.Duration) bool {
	select {
	case GetNotificationPort(id) <- data:
		return true
	case <-time.After(timeout):
		return false
	}
}

func NotifyNoWait(id string, data interface{}) bool {
	select {
	case GetNotificationPort(id) <- data:
		return true
	default:
		return false
	}
}
