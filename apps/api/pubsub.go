package main

import "sync"

type PubsubMsg struct {
	Topic   string
	Message []byte
}

type ISubscriber interface {
	Subscribe(topic string)
	Unsubscribe(topic string)
	GetChannel() chan PubsubMsg
	Close()
}

type IPubsub interface {
	Send(topic string, msg []byte) error
	AddSubscriber() ISubscriber
}

type Subscriber struct {
	mu      sync.RWMutex
	Topics  []string
	Channel chan PubsubMsg
	pubsub  *Pubsub
}

func (sub *Subscriber) Subscribe(topic string) {
	sub.mu.Lock()
	defer sub.mu.Unlock()
	for _, tpc := range sub.Topics {
		if tpc == topic {
			return
		}
	}
	sub.Topics = append(sub.Topics, topic)
}

func (sub *Subscriber) Unsubscribe(topic string) {
	sub.mu.Lock()
	defer sub.mu.Unlock()
	index := findIndex(sub.Topics, topic)
	if index != -1 {
		sub.Topics = remove(sub.Topics, index)
	}
}

func (sub *Subscriber) Close() {
	sub.pubsub.mu.Lock()
	sub.mu.Lock()
	defer sub.pubsub.mu.Unlock()
	defer sub.mu.Unlock()
	sub.Topics = []string{}
	close(sub.Channel)
	index := findIndex(sub.pubsub.Subscribers, sub)
	if index != -1 {
		remove(sub.pubsub.Subscribers, index)
	}
}

func (sub *Subscriber) GetChannel() chan PubsubMsg {
	return sub.Channel
}

type Pubsub struct {
	mu          sync.RWMutex
	Subscribers []*Subscriber
}

var PubsubClient IPubsub = &Pubsub{}

func (ps *Pubsub) Send(topic string, msg []byte) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	payload := PubsubMsg{topic, msg}
	for _, sub := range ps.Subscribers {
		sub.mu.RLock()
		index := findIndex(sub.Topics, topic)
		if index == -1 {
			continue
		}
		sub.Channel <- payload
		sub.mu.RUnlock()
	}
	return nil
}
func (ps *Pubsub) AddSubscriber() ISubscriber {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	sub := &Subscriber{Channel: make(chan PubsubMsg, 10), pubsub: ps}
	ps.Subscribers = append(ps.Subscribers, sub)
	return sub
}
