package pubsub

import (
	"fmt"
	"main/websnake"
	"net"
	"sync"
)

var globalPubSubService *PubSubService
var once sync.Once

// GetGlobalPubSubService возвращает глобальный экземпляр сервиса Pub/Sub
func GetGlobalPubSubService() *PubSubService {
	once.Do(func() {
		globalPubSubService = NewPubSubService()
	})
	return globalPubSubService
}

// PubSubService представляет сервис Pub/Sub
type PubSubService struct {
	subscribers map[string][]Subscriber
	mu          sync.Mutex
}

// NewPubSubService создает новый экземпляр сервиса Pub/Sub
func NewPubSubService() *PubSubService {
	return &PubSubService{
		subscribers: make(map[string][]Subscriber),
	}
}

// Subscribe подписывает канал на указанную тему
func (p *PubSubService) Subscribe(topic string, sub Subscriber) {
	fmt.Printf("%+v", p.subscribers)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers[topic] = append(p.subscribers[topic], sub)
	fmt.Println("Subscribe for topic ", topic)
	fmt.Printf("%+v", p.subscribers)
}

// Publish отправляет сообщение в указанную тему
func (p *PubSubService) Publish(topic string, msg Message) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if subscribers, ok := p.subscribers[topic]; ok {
		for _, sub := range subscribers {
			go func(msg Message) {
				sub.EventHandler(msg)
			}(msg)
		}
	}
}

// Subscriber представляет собой подписчика
type Subscriber struct {
	EventChannel chan string
	EventHandler func(msg Message)
}

type Message struct {
	Msg  *websnake.GameMessage
	From *net.UDPAddr
	To   *net.UDPAddr
}
