package main

import (
	"fmt"
	"os"

	"github.com/goiiot/libmqtt"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s [server:port] [topic] [message]\n", os.Args[0])
		os.Exit(1)
	}

	hostname, _ := os.Hostname()
	clientID := fmt.Sprintf("%s-on-%s", os.Args[0], hostname)
	client, err := libmqtt.NewClient(
		libmqtt.WithServer(os.Args[1]),
		libmqtt.WithClientID(clientID),
	)

	if err != nil {
		// handle client creation error
		fmt.Println("Can't create client", err)
	}
	// connect to server

	// success
	// you are now connected to the `server`
	// (the `server` is one of your provided `servers` when create the client)
	// start your business logic here or send a signal to your logic to start

	// register publish handler
	{
		client.HandleNet(func(server string, err error) {
			if err != nil {
				fmt.Printf("error happened to connection to server [%v]: %v\n", server, err)
			}
		})

		client.HandlePub(func(topic string, err error) {
			if err != nil {
				fmt.Printf("publish packet to topic [%v] failed: %v\n", topic, err)
			} else {
				fmt.Printf("publish [%v] to topic [%v] on [%v]: ok\n", os.Args[3], topic, os.Args[1])
			}
			client.Destroy(true)
		})
	}

	// publish some topic message(s)
	// connect to server
	client.Connect(func(server string, code byte, err error) {
		if err != nil {
			fmt.Printf("connect to server [%v] failed: %v\n", server, err)
			return
		}

		if code != libmqtt.CodeSuccess {
			fmt.Printf("connect to server [%v] failed with server code [%v]\n", server, code)
			return
		}

		// connected
		client.Publish(
			&libmqtt.PublishPacket{
				TopicName: os.Args[2],
				Qos:       libmqtt.Qos0,
				Payload:   []byte(os.Args[3]),
			})
	})
	client.Wait()

}
