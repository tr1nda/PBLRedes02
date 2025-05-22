package mqtt

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TopicoHealthCheck       = "healthcheck/servers"
	TopicoStatus            = "servers/status"
	TopicoReservaConfirmada = "reservas/confirmadas"
	TopicoReservaNegada     = "reservas/negadas"
	TopicoPontosDisponiveis = "pontos/disponiveis"
)

var Broker = ""

func StartBroker() {
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://mqtt:1883"
	}

	Broker = broker
}

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
