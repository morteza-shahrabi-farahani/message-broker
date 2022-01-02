package broker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"sync"
	"therealbroker/pkg/broker"
	"time"
	//log "github.com/sirupsen/logrus"
)

const (
	host     = "postgres"
	dbport   = "5432"
	user     = "postgres"
	password = "12345678"
	dbname   = "broker"
)

var IDCounter = 0
var publishStatement = `INSERT INTO MESSAGES2 (subject, body,startTime, expirationTime) VALUES ($1, $2, $3, $4) RETURNING id;`
var deleteStatement = `DELETE FROM MESSAGES2 WHERE id = $1;`
var selectStatement = `SELECT * FROM users WHERE id=$1;`
var dbUrl = "postgres://" + user + ":" + password + "@" + host + ":" + dbport + "/?pool_max_conns=50"

type Module struct {
	isClosed    bool
	isConnected bool
	chats       []chat
	counter     int
	allMessages []messages
	dbConnect   *pgxpool.Pool
}

type messages struct {
	payam      broker.Message
	expireTime time.Duration
	startTime  time.Time
	id         int
}

type chat struct {
	name               string
	messages           []broker.Message
	subscribedChannels []chan broker.Message
}

func NewModule() broker.Broker {
	return &Module{}
}

func (m *Module) Close() error {
	var module Module
	module.isClosed = true
	return nil
}

func createMessageNew(body string, expirationTime float64) broker.Message {

	return broker.Message{
		Body:       string(body),
		Expiration: time.Duration(expirationTime * 1000000000),
	}
}

func (m *Module) ConnectDatabase() error {

	dbConn, err9 := pgxpool.Connect(context.Background(), dbUrl)
	if err9 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err9)
	}
	m.dbConnect = dbConn
	m.isConnected = true

	sqlStatement1 := "CREATE TABLE IF NOT EXISTS MESSAGES2 (id SERIAL PRIMARY KEY, subject VARCHAR(50), body VARCHAR(500), startTime TIME, expirationTime FLOAT);"
	_, err8 := m.dbConnect.Exec(context.Background(), sqlStatement1)
	if err8 != nil {
		fmt.Println("error for creating tablle.\n")
		return errors.New("couldn't create table")
	}

	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.isClosed == true {
		fmt.Println("this flag is false", m.isClosed)
		return 0, errors.New("service is unavailable")
	}

	if !m.isConnected {
		m.ConnectDatabase()
	}
	var lock sync.Mutex
	var lock3 sync.Mutex
	//m.CreateTable()
	lock.Lock()
	var hasChat = false
	//m.counter++
	err6 := m.dbConnect.QueryRow(context.Background(), publishStatement, subject, msg.Body, time.Now(), msg.Expiration).Scan(&m.counter)
	if err6 != nil {
		log.Printf("unable to insert to database because %v", err6)
	} /*else {
		log.Printf("id is %v", m.counter)
	}*/

	for i := 0; i < len(m.chats); i++ {
		if m.chats[i].name == subject {
			lock3.Lock()
			var lock2 sync.Mutex
			hasChat = true
			//fmt.Println("hello")
			m.chats[i].messages = append(m.chats[i].messages, msg)
			var payam messages
			payam.payam = msg
			payam.id = m.counter
			payam.expireTime = msg.Expiration
			payam.startTime = time.Now()
			m.allMessages = append(m.allMessages, payam)
			//lock.Unlock()
			for j := 0; j < len(m.chats[i].subscribedChannels); j++ {
				//fmt.Println(len(m.chats[i].subscribedChannels))
				lock2.Lock()
				m.chats[i].subscribedChannels[j] <- msg
				lock2.Unlock()
			}
			lock3.Unlock()
			break
		}
	}

	if hasChat == false {
		lock3.Lock()
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
		lock3.Unlock()
	}

	lock.Unlock()
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
		return broker.Message{}, errors.New("service is unavailable")
	}

	var selectId int
	var selectSubject string
	var selectBody string
	var selectStartTime time.Time
	var selectExpirationTime float64

	row := m.dbConnect.QueryRow(context.Background(), selectStatement, id)
	switch err := row.Scan(&selectId, &selectSubject, &selectBody, &selectStartTime, &selectExpirationTime); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		if time.Now().Sub(selectStartTime) < time.Duration(selectExpirationTime) {
			message := createMessageNew(selectBody, selectExpirationTime)
			return message, nil
		} else {
			return broker.Message{}, errors.New("message with id provided is expired")
		}

	}
	//var lock sync.Mutex
	for i := 0; i < len(m.allMessages); i++ {
		if m.allMessages[i].id == id {
			//lock.Lock()
			if time.Now().Sub(m.allMessages[i].startTime) < m.allMessages[i].expireTime {
				return m.allMessages[i].payam, nil
			} else {
				return broker.Message{}, errors.New("message with id provided is expired")
			}
			//lock.Unlock()
		}
	}

	return broker.Message{}, nil
}
