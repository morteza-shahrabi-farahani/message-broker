package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"therealbroker/api/proto"
)

const (
	address = "localhost:5050"
)

//var pingCounter2 = prometheus.NewCounter(
//	prometheus.CounterOpts{
//		Name: "ping_request_count",
//		Help: "No of request handled by Ping handler",
//	},
//)
//
//func ping2(w http.ResponseWriter, req *http.Request) {
//	pingCounter.Inc()
//	fmt.Fprintf(w, "pong")
//}

func main() {

	//go func() {
	//	//router := mux.NewRouter()
	//	//router.Use(prometheusMiddleware)
	//	//router.Path("/prometheus").Handler(promhttp.Handler())
	//
	//	// Serving static files
	//	//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//
	//	//fmt.Println("Serving requests on port 9000")
	//	//err := http.ListenAndServe(":9000", router)
	//	//log.Fatal(err)
	//	//prometheus.MustRegister(pingCounter2)
	//	//http.Handle("/metrics", promhttp.Handler())
	//	//http.HandleFunc("/ping", ping2)
	//	//http.ListenAndServe(":5052", nil)
	//
	//}()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewBrokerClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	publishTest := proto.PublishRequest{
		Subject:           "ali",
		Body:              []byte("110 110"),
		ExpirationSeconds: 100000000,
	}
	fmt.Println(publishTest)
	publishResponse, newErr := c.Publish(context.Background(), &publishTest)
	if newErr != nil {
		log.Fatalf("error is: %v", newErr)
	}
	log.Printf("detail of response. id is %v", publishResponse.GetId())
	fetchTest := proto.FetchRequest{
		Subject: "ali",
		Id:      publishResponse.GetId(),
	}
	body, err2 := c.Fetch(context.Background(), &fetchTest)
	if err2 != nil {
		log.Fatalf("error is %v", err2)
	}
	log.Printf("detail of fetch request is %v", body)
}
