package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/srthkh2k/grpc-go/chatserver"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Enter Server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connecting : " + serverID)

	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close()

	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	ch := clienthandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	bl := make(chan bool)
	<-bl

}

type clienthandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clienthandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")

}

func (ch *clienthandle) sendMessage() {

	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}

func (ch *clienthandle) receiveMessage() {

	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)

	}
}
