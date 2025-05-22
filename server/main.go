package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pblredes2/pkg/shared/mqtt"
	server "pblredes2/server/internal/handler"
	"pblredes2/server/internal/model"
)

var (
	Pontos = []model.PontoRecarga{}
)

func main() {
	err := model.CarregarPontos()
	if err != nil {
		log.Fatalf("Erro ao carregar pontos de recarga: %v", err)
	}

	// Debug: imprime os pontos carregados
	for _, p := range model.Pontos {
		fmt.Printf("Ponto: %s | Região: %s | Fila: %v\n", p.ID, p.Regiao, p.Fila)
	}

	http.HandleFunc("/", server.Handler)
	http.HandleFunc("/iniciar_rota", server.IniciarRota)
	http.HandleFunc("/verificar_pontos", server.ListarPontos)

	porta := ":9000"
	fmt.Printf("Servidor rodando na porta %s...\n", porta)

	mqtt.StartBroker()

	serverId := fmt.Sprintf("server-%s", os.Getenv("INSTANCE_ID"))
	client := mqtt.Connect(mqtt.Broker, serverId)

	mqtt.Subscribe(client, mqtt.TopicoStatus, server.DefaultHandler)
	mqtt.Subscribe(client, mqtt.TopicoHealthCheck, server.HandleStatusServidores)
	mqtt.Subscribe(client, mqtt.TopicoPontosDisponiveis, server.HandlePontosDisponiveis)

	// Inicia o servidor na porta 9000
	err = http.ListenAndServe(porta, nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
