package server

import (
	"fmt"
	"log"
	"os"

	"pblredes2/pkg/shared/mqtt"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
)

func DefaultHandler(client mqttlib.Client, message mqttlib.Message) {
	receivedMessage(message)
}

// Handle para verificar o status dos servidores
func HandleStatusServidores(client mqttlib.Client, message mqttlib.Message) {
	receivedMessage(message)
	mqtt.Publish(client, mqtt.TopicoStatus, fmt.Sprintf("{\"id\": \"%s\",\"status\": %t", serverID(), true))
}

// Handle para retorno de pontos disponíveis no servidor
func HandlePontosDisponiveis(client mqttlib.Client, message mqttlib.Message) {
	receivedMessage(message)
	mqtt.Publish(client, mqtt.TopicoStatus, fmt.Sprintf("{\"id\": \"%s\",\"status\": %t", serverID(), true))
}

func serverID() string {
	return fmt.Sprintf("server-%s", os.Getenv("INSTANCE_ID"))
}

func receivedMessage(message mqttlib.Message) {
	log.Printf("Mensagem recebida no tópico %s:\n%s", message.Topic(), string(message.Payload()))
}
