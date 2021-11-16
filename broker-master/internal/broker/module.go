package broker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"therealbroker/pkg/broker"
	"time"
)

type Module struct {
	isClosed bool
	chats []chat
	counter int
	allMessages []messages
}

type messages struct {
	payam broker.Message
	expireTime time.Duration
	startTime time.Time
	id int
}

type chat struct {
	name string
	messages []broker.Message
	subscribedChannels []chan broker.Message
}

//var subscribedChannelsList []chan broker.Message
type subscribedChats struct {

}

func NewModule() broker.Broker {
	return &Module{}
}

func (m *Module) Close() error {
	var module Module
	module.isClosed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.isClosed == true {
		fmt.Println("this flag is false", m.isClosed)
		return 0, errors.New("service is unavailable")
	}

	var lock sync.Mutex
	m.counter++
	var hasChat = false
	for i := 0; i < len(m.chats); i++ {
		if m.chats[i].name == subject {
			lock.Lock()
			hasChat = true
			//fmt.Println("hello")
			m.chats[i].messages = append(m.chats[i].messages, msg)
			var payam messages
			payam.payam = msg
			payam.id = m.counter
			payam.expireTime = msg.Expiration
			payam.startTime = time.Now()
			m.allMessages = append(m.allMessages, payam)
			for j := 0; j < len(m.chats[i].subscribedChannels); j++ {
				//fmt.Println(len(m.chats[i].subscribedChannels))
				m.chats[i].subscribedChannels[j] <- msg
			}
			lock.Unlock()
			break
		}
	}

	if hasChat == false {
		lock.Lock()
		var newChat chat
		newChat.name = subject
		m.chats = append(m.chats, newChat)
		newChat.messages = append(newChat.messages, msg)
		var payam messages
		payam.payam = msg
		payam.id = m.counter
		payam.expireTime = msg.Expiration
		payam.startTime = time.Now()
		m.allMessages = append(m.allMessages, payam)
		lock.Unlock()
	}

	return m.counter, nil

}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.isClosed == true {
		return nil, errors.New("service is unavailable")
	}

	var hasChat = false
	//var wg sync.WaitGroup
	var lock sync.Mutex
	newSubChan := make(chan broker.Message, 100)
	for i := 0; i < len(m.chats); i++ {
		if m.chats[i].name == subject {
			lock.Lock()
			hasChat = true
			m.chats[i].subscribedChannels = append(m.chats[i].subscribedChannels, newSubChan)
			lock.Unlock()
		}
	}


	if hasChat == false {
		lock.Lock()
		var newChat chat
		newChat.name = subject
		newChat.subscribedChannels = append(newChat.subscribedChannels, newSubChan)
		m.chats = append(m.chats, newChat)
		lock.Unlock()
	}

	return newSubChan, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	if m.isClosed == true {
		return broker.Message{

		}, errors.New("service is unavailable")
	}

	for i := 0; i < len(m.allMessages); i++ {
		if m.allMessages[i].id == id {
			if time.Now().Sub(m.allMessages[i].startTime) < m.allMessages[i].expireTime {
				return m.allMessages[i].payam, nil
			} else {
				return broker.Message{}, errors.New("message with id provided is expired")
			}
		}
	}

	return broker.Message{}, nil
}

