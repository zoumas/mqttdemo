package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Replace with your MQTT broker address
	const broker = "tcp://test.mosquitto.org:1883"
	client := mqtt.NewClient(configuredMQTTClientOptions(broker))

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Subscribe to a topic
	const topic = "your/topic"
	if token := client.Subscribe(topic, 1, func(_ mqtt.Client, m mqtt.Message) {
		fmt.Printf("Received message on topic %q: %q\n", m.Topic(), m.Payload())
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Publish a message
	const message = "Hello, MQTT!"
	client.Publish(topic, 1, false, message).Wait()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	client.Disconnect(250)
}

func configuredMQTTClientOptions(brokerAddr string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(brokerAddr)
	opts.SetClientID("example-client")

	// Doesn't seem to work 😭
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("Connected", c.IsConnected())
	})
	return opts
}
