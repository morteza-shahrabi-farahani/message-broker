package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
	//"github.com/gorilla/mux"
	proto2 "therealbroker/api/proto"
	broker2 "therealbroker/internal/broker"
	broker3 "therealbroker/pkg/broker"
	"time"
)

const (
	port = ":5050"
)

type BrokerServer struct {
	proto2.UnimplementedBrokerServer
	module broker2.Module
}

func createMessageNew(request *proto2.PublishRequest) broker3.Message {

	return broker3.Message{
		Body:       string(request.GetBody()),
		Expiration: time.Duration(request.GetExpirationSeconds() * 1000000000),
	}
}

func (server *BrokerServer) Publish(ctx context.Context, publishRequest *proto2.PublishRequest) (publishResponse *proto2.PublishResponse, err error) {
	timer := prometheus.NewTimer(durationHistogram.WithLabelValues("publish duration"))
	//startTime := time.Now()
	var lock sync.Mutex
	var result proto2.PublishResponse
	msg := createMessageNew(publishRequest)
	fmt.Println(msg)
	var msgId int
	lock.Lock()
	msgId, err = server.module.Publish(ctx, publishRequest.GetSubject(), msg)
	if err != nil {
		//duration := time.Since(startTime)
		timer.ObserveDuration()
		totalRequests.WithLabelValues("failed").Inc()
		//durationHistogram.WithLabelValues("publish latecy").Observe(float64(duration))
		return nil, err
	}

	result.Id = int32(msgId)
	lock.Unlock()
	totalRequests.WithLabelValues("succeed").Inc()
	timer.ObserveDuration()
	return &result, nil
}

func (server *BrokerServer) Subscribe(subscribeRequest *proto2.SubscribeRequest, subscriveServer proto2.Broker_SubscribeServer) error {
	timer := prometheus.NewTimer(durationHistogram.WithLabelValues("subscribe duration"))
	ch, err := server.module.Subscribe(context.Background(), subscribeRequest.GetSubject())
	if err != nil {
		timer.ObserveDuration()
		totalRequests.WithLabelValues("failed").Inc()
		return err
	} else {
		totalRequests.WithLabelValues("succeed").Inc()
		totalActiveSubscribed.WithLabelValues("subscribed channels").Inc()
		for i := 0; i < len(ch); i++ {
			msg := <-ch
			msgResponse := proto2.MessageResponse{Body: []byte(msg.Body)}
			subscriveServer.Send(&msgResponse)
		}
		timer.ObserveDuration()
		return nil
	}

}

func (server *BrokerServer) Fetch(ctx context.Context, fetchRequest *proto2.FetchRequest) (*proto2.MessageResponse, error) {
	timer := prometheus.NewTimer(durationHistogram.WithLabelValues("fetch duration"))
	var messageResponse proto2.MessageResponse
	msg, err := server.module.Fetch(ctx, fetchRequest.Subject, int(fetchRequest.Id))
	if err != nil {
		timer.ObserveDuration()
		totalRequests.WithLabelValues("failed").Inc()
		return nil, err
	} else {
		timer.ObserveDuration()
		totalRequests.WithLabelValues("succeed").Inc()
		messageResponse.Body = []byte(msg.Body)
		return &messageResponse, nil
	}
}

//metrics
var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "total_requests",
		Help: "Number of requests.",
	},
	[]string{"statusCode"},
)

var durationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "call_duration",
	Help: "Duration of calls.",
}, []string{"path"})

var totalActiveSubscribed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "number_of_subscribed_channels",
		Help: "number of subsctrived channels",
	},
	[]string{"path"},
)

//func prometheusMiddleware(next http.Handler) http.Handler {
// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	 rw := NewResponseWriter(w)
//	 next.ServeHTTP(rw, r)
//
//	 totalRequests.WithLabelValues(path).Inc()
// })
//}

//func init() {
// prometheus.Register(totalRequests)
//}

func main() {

	go func() {
		//router := mux.NewRouter()
		//router.Use(prometheusMiddleware)
		//router.Path("/prometheus").Handler(promhttp.Handler())

		// Serving static files
		//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

		//fmt.Println("Serving requests on port 9000")
		//err := http.ListenAndServe(":9000", router)
		//log.Fatal(err)
		prometheus.Register(totalRequests)
		prometheus.Register(durationHistogram)
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":5051", nil)

	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto2.RegisterBrokerServer(s, &BrokerServer{})
	log.Printf("server listening at: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
