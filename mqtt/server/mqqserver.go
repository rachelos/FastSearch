package server

import (
	"fmt"
	"log"
	"os"

	mqtt "gitee.com/rachel_os/fastsearch/mochi-mqtt/server"
	"gitee.com/rachel_os/fastsearch/mochi-mqtt/server/hooks/auth"
	"gitee.com/rachel_os/fastsearch/mochi-mqtt/server/packets"
	"gopkg.in/yaml.v2"
)

type MQTTServer struct {
	Server *mqtt.Server
	Path   string
}

func NewMQTTServer(path string) *MQTTServer {
	return &MQTTServer{
		Path: path,
	}
}
func (s *MQTTServer) Run() {
	// 创建信号用于等待服务端关闭信号
	// sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	<-sigs
	// 	done <- true
	// }()

	// 创建新的 MQTT 服务器。
	var config *mqtt.Options
	file, err := os.ReadFile(s.Path) //详情：https://github.com/golang/go/issues/42026
	if err != nil {
		panic(err)
	}
	_ = yaml.Unmarshal(file, &config)
	server := mqtt.New(config)

	// 允许所有连接(权限)。
	_ = server.AddHook(new(auth.AllowHook), nil)

	// 在标1883端口上创建一个 TCP 服务端。
	// tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: ":1883"})
	// err := server.AddListener(tcp)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	s.Server = server
	// 服务端等待关闭信号
	// <-done

	// 关闭服务端时需要做的一些清理工作
}
func (s *MQTTServer) Hook(topic string, callbackFn func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet)) {
	if callbackFn == nil {
		callbackFn = func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
			fmt.Println("received message", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
		}
	}

	s.Server.Subscribe(topic+"/#", 1, callbackFn)
}
func (s *MQTTServer) Publish(topic string, message []byte) error {
	err := s.Server.Publish(topic, message, false, 0)
	return err
}
