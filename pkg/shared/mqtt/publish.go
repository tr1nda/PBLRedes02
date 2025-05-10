package mqtt

import (
	"log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Publish envia uma mensagem para o tópico especificado
func Publish(client mqtt.Client, topic string, payload interface{}) {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()

	if err := token.Error(); err != nil {
		log.Printf("Erro ao publicar no tópico %s: %v", topic, err)
	}
}
