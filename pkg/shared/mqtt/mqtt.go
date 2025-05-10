package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

// Connect cria um novo cliente MQTT e conecta ao broker
func Connect(brokerURL, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.AutoReconnect = true

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		log.Fatalf("Erro ao conectar ao MQTT: %v", err)
	}

	log.Println("Conectado ao broker MQTT:", brokerURL)
	return client
}
