package main

import (
	"fmt"
	"net/http"
	"os"
	"pblredes2/pkg/shared/mqtt"
	server "pblredes2/server/internal/handler"
)

const (
	topicoTeste             = "topico/teste"
	topicoStatus            = "servers/status"
	topicoReservaConfirmada = "reservas/confirmadas"
	topicoReservaNegada     = "reservas/negadas"
	topicoPontosDisponiveis = "pontos/disponiveis"
)

func main() {
	http.HandleFunc("/", server.Handler)
	http.HandleFunc("/iniciar_rota", server.IniciarRota) // Define a função que vai tratar as requisições para "/"

	porta := ":9000"
	fmt.Printf("Servidor rodando na porta %s...\n", porta)

	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://mqtt:1883"
	}

	client := mqtt.Connect(broker, "server-id")

	mqtt.Subscribe(client, topicoTeste, mqtt.DefaultHandler)
	mqtt.Subscribe(client, topicoStatus, mqtt.DefaultHandler)

	mqtt.Publish(client, topicoTeste, "Olá do container!")
	mqtt.Publish(client, topicoTeste, "Nova mensagem no broker!")

	estaDisponivel := true
	go func(disponivel bool) {
		fmt.Printf("Disponível: %t\n", disponivel)
		for i := 0; i < 5; i++ {
			mqtt.Publish(client, topicoStatus, fmt.Sprintf("{\"id\": \"%s\",\"status\": %t", "server1", disponivel))
		}
	}(estaDisponivel)

	// Inicia o servidor na porta 9000
	err := http.ListenAndServe(porta, nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
