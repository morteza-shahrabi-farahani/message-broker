package main

import (
	"fmt"
	"math/rand"
	"therealbroker/pkg/broker"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

//type server struct {}

func main() {
	fmt.Println("Hello!")
	// listener, err := net.Listen("tcp", ":4040")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// srv := grpc.NewServer()
	// if err2 := srv.Serve(listener); err2 != nil {
	// 	log.Fatalf("failed to serve: %s", err2)
	// }

	// //service broker.Broker
	// //service = broker2.NewModule()
	// //mainCtx = context.Background()
	// brokHolder := broker2.Module{}
	// msg := createMessage()
	// brokHolder.Publish(context.Background(), "", msg)
	//broker.Broker().Publish()
	//publishTest := proto.PublishRequest{
	//	Subject:           "ali",
	//	Body: []byte("nil"),
	//	ExpirationSeconds: 10,
	//}
	//fmt.Println(publishTest)

	//func (m *server) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	//	if m.isClosed == true {
	//		fmt.Println("this flag is false", m.isClosed)
	//		return 0, errors.New("service is unavailable")
	//	}
	//}

	//data, err := proto.Marshal(publishTest)
	//fmt.Println(data)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createMessage() broker.Message {
	body := randomString(16)

	return broker.Message{
		Body:       body,
		Expiration: 0,
	}
}
