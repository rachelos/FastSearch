package client

import (
	"fmt"
	"log"
	"testing"

	mqtt1 "gitee.com/rachel_os/fastsearch/mochi-mqtt/server"
	"gitee.com/rachel_os/fastsearch/mochi-mqtt/server/packets"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestClient(t *testing.T) {
	server := mqtt1.New(&mqtt1.Options{
		InlineClient: true,
	})
	err := server.Publish("direct/publish", []byte("packet scheduled message"), false, 0)
	fmt.Println(err)
	callbackFn := func(cl *mqtt1.Client, sub packets.Subscription, pk packets.Packet) {
		server.Log.Info("inline client received message from subscription", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
	}
	server.Subscribe("direct/#", 1, callbackFn)

	server.Unsubscribe("direct/#", 1)
}

func Test2(t *testing.T) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("fastsearch1")
	opts.SetAutoReconnect(true)

	opts.SetUsername("a") // 如果需要的话
	opts.SetPassword("b") // 如果需要的话

	opts.SetOnConnectHandler(mqtt.OnConnectHandler(onConnect))
	opts.SetConnectionLostHandler(mqtt.ConnectionLostHandler(onConnectionLost))
	opts.SetDefaultPublishHandler(mqtt.MessageHandler(onMessage))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT server: %v", token.Error())
	}

	// 订阅一个或多个主题
	if token := client.Subscribe("direct/#", 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic: %v", token.Error())
	}
	client.Connect().Done()
	// 阻塞主goroutine，直到收到退出信号或其他条件
	select {}
}

func onConnect(client mqtt.Client) {
	fmt.Println("Connected to MQTT server")
}

func onConnectionLost(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func onMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", msg.Topic(), msg.Payload())
}
