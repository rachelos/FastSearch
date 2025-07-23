package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/rachel_os/fastsearch/global"
	mqtt "gitee.com/rachel_os/fastsearch/mochi-mqtt/server"
	"gitee.com/rachel_os/fastsearch/mochi-mqtt/server/packets"
	"gitee.com/rachel_os/fastsearch/searcher/model"
)

var MQTT_SERVER *MQTTServer

func (s *MQTTServer) InitRouter(topic string) {
	MQTT_SERVER = s
	s.Hook(topic, func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		go func() {
			switch strings.ReplaceAll(pk.TopicName, topic, "") {
			case "/index":
				document := &model.IndexDoc{}
				json.Unmarshal([]byte(string(pk.Payload)), &document)
				err := global.Container.GetDataBase("default").IndexDocument(document)
				if err != nil {
					fmt.Println(err)
					s.Publish(pk.TopicName+"/err", []byte(err.Error()))
				} else {
					s.Publish(pk.TopicName+"/back", []byte(document.Id))
					// fmt.Println("received message", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
				}
			}
		}()
	})
}
