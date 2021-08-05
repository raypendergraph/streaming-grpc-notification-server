package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"server/generated/notifications"
	"time"
)
const connections = 250000

func main()  {
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := notifications.NewNotificationServiceClient(conn)
	for i := 0; i < connections; i+=1 {
		go runTest(i, client)
	}
	var done chan interface{}
	<- done
}

func runTest(id int, client notifications.NotificationServiceClient) {
	idStr := fmt.Sprintf("%d", id)
	req := notifications.SubscribeRequest{Id: idStr}
	stream, err := client.Subscribe(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	<-time.NewTimer(1 * time.Second).C
	for {
		//send random
		sendTestMessage(idStr, client)

		//block receive
		_, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		//What a moment
		<-time.NewTimer(time.Duration(rand.Intn(1500))  * time.Millisecond).C
	}

}

func sendTestMessage(senderStr string, client notifications.NotificationServiceClient) {
	recipientStr := fmt.Sprintf("%d", rand.Intn(connections))
	_, err := client.SendMessage(context.Background(), &notifications.SendMessageRequest{
		Sender: senderStr,
		Recipient: recipientStr,
		Message:   fmt.Sprintf("%s->%s", senderStr, recipientStr),
	})
	if err != nil {
		panic(err)
	}
}