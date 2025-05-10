package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pblredes2/pkg/shared/mqtt"
	server "pblredes2/server/internal/handler"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	http.HandleFunc("/", server.Handler) // Define a função que vai tratar as requisições para "/"

	porta := ":8080"
	fmt.Printf("Servidor rodando na porta %s...\n", porta)

	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://mqtt:1883"
	}

	client := mqtt.Connect(broker, "server-id")

	mqtt.Subscribe(client, "topico/teste", func(c mqttlib.Client, m mqttlib.Message) {
		log.Printf("Recebido: %s", string(m.Payload()))
	})

	mqtt.Publish(client, "topico/teste", "Olá do container!")

	// Inicia o servidor na porta 8080
	err := http.ListenAndServe(porta, nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}

	select {} // impede o encerramento
}
