package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

type PontoRecarga struct {
	ID     string
	Regiao string
	Fila   []string
}

// Função para carregar os jsons da pasta data
func CarregarPontos(path string) ([]PontoRecarga, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pontos []PontoRecarga
	err = json.Unmarshal(data, &pontos)
	if err != nil {
		return nil, err
	}

	return pontos, nil
}

func main() {
	// Lê nome do arquivo da variável de ambiente
	jsonFile := os.Getenv("JSON_FILE")
	if jsonFile == "" {
		log.Fatal("Variável de ambiente JSON_FILE não foi definida")
	}

	// Caminho dentro do container
	fullPath := filepath.Join("data", jsonFile)

	pontos, err := CarregarPontos(fullPath)
	if err != nil {
		log.Fatalf("Erro ao carregar pontos de recarga: %v", err)
	}

	// Debug: imprime os pontos carregados
	for _, p := range pontos {
		fmt.Printf("Ponto: %s | Região: %s | Fila: %v\n", p.ID, p.Regiao, p.Fila)
	}

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
	err2 := http.ListenAndServe(porta, nil)
	if err2 != nil {
		fmt.Println("Erro ao iniciar o servidor:", err2)
	}
}
