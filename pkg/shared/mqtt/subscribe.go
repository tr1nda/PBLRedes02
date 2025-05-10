package mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Subscribe registra um callback para um tópico
func Subscribe(client mqtt.Client, topic string, handler mqtt.MessageHandler) {
	token := client.Subscribe(topic, 0, handler)
	token.Wait()

	if err := token.Error(); err != nil {
		log.Fatalf("Erro ao assinar tópico %s: %v", topic, err)
	} else {
		log.Println("Inscrito no tópico:", topic)
	}
}
