package main

import (
	"context"
	"fmt"
	"github.com/pkg/profile"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "server/generated/notifications"
	"sync"
)

var (
	clients       = make(map[string]chan pb.MessageEvent)
	mutex         sync.RWMutex
	cancellations = make(chan string)
)

func getAllClients() []chan<-pb.MessageEvent {
	mutex.RLock()
	defer mutex.RUnlock()
	response := make([]chan<-pb.MessageEvent, len(clients))
	i := 0
	for _, v := range clients{
		response[i] = v
		i+=1
	}
	return response
}

func getClient(key string) chan<-pb.MessageEvent {
	mutex.RLock()
	defer mutex.RUnlock()
	return clients[key]
}

func removeClient(key string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clients, key)
}

func addClient(key string, c chan pb.MessageEvent) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := clients[key]; ok {
		return fmt.Errorf("%s is already registered", key)
	}
	clients[key] = c
	return nil
}

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
}

func (r NotificationHandler) Subscribe(request *pb.SubscribeRequest, server pb.NotificationService_SubscribeServer) error {
	fmt.Printf("Requesting subscription for %s\n", request.Id)
	clientMessages := make(chan pb.MessageEvent, 8)
	addClient(request.Id, clientMessages)
	defer removeClient(request.Id)
	return handle(request.Id, server, clientMessages, server.Context().Done())
}

func (r NotificationHandler) Unsubscribe(ctx context.Context, req *pb.UnsubscribeRequest) (*pb.UnsubscribeResponse, error) {
	fmt.Printf("Requesting unsubscribe for %s\n", req.Id)
	removeClient(req.Id)
	return &pb.UnsubscribeResponse{}, nil
}

func (r NotificationHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	event := pb.MessageEvent{
		Sender:      req.Sender,
		Message:     req.Message,
	}
	if req.Recipient == "" {
		event.IsBroadcast = true
		//broadcast
		//fmt.Printf("broadcasting message: %s\n", req.Message)
		for _, c := range getAllClients() {
			c <- event
		}
	} else {
		//fmt.Printf("sending message to [%s]: %s\n", req.Recipient, req.Message)
		if c := getClient(req.Recipient); c != nil{
			c <- event
		}
	}
	return &pb.SendMessageResponse{}, nil
}

func main() {
	defer profile.Start(
		profile.ProfilePath("./profiles"),
		//profile.MemProfile,
		//profile.CPUProfile,
		profile.GoroutineProfile,
		).Stop()

	done := make(chan interface{})
	server := grpc.NewServer()
	pb.RegisterNotificationServiceServer(server, NotificationHandler{})

	go func() {
		defer close(done)
		listen, err := net.Listen("tcp", ":8888")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		fmt.Printf("grpc server starting on [%s] \n", listen.Addr().String())
		panic(server.Serve(listen))
	}()
	<-done
	fmt.Println("It's been fun...")
}

func handle(id string, server pb.NotificationService_SubscribeServer, messages <-chan pb.MessageEvent, done <-chan struct{}) error {
	fmt.Println("starting a mux")
	for {
		select {
		case message, ok := <-messages:
			if !ok {
				fmt.Printf("%s message channel was closed.\n", id)
				return nil
			}
			//fmt.Printf("%s: sending a message\n", id)
			err := server.Send(&message)
			if err != nil {
				return err
			}
		case <-done:
			fmt.Printf("%s: server instance is telling us to bail. Bye.\n", id)
			return nil
		}
	}
}
