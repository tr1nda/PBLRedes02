package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"pblredes2/pkg/shared/mqtt"
	"pblredes2/server/internal/model"
	"strings"
)

// var serversSlice = []string{"http://servidor1:9000", "http://servidor2:9000", "http://servidor3:9000"}

func nextServer(currentID string) string {
	for i, s := range servers {
		if strings.Contains(s, currentID) {
			return servers[(i+1)%len(servers)]
		}
	}
	return ""
}

const (
	Server1 = "http://servidor1:9000/"
	Server2 = "http://servidor2:9000/"
	Server3 = "http://servidor3:9000/"
)

var (
	servers = map[int]string{
		1: Server1,
		2: Server2,
		3: Server3,
	}
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Isso é basicamente um health check!")

	serverId := fmt.Sprintf("server-%s", os.Getenv("INSTANCE_ID"))
	client := mqtt.Connect(mqtt.Broker, serverId)
	mqtt.Publish(client, mqtt.TopicoHealthCheck, "ping")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tá tudo funcionando!"))
}

func IniciarRota(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Iniciando rota...")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var reserva model.Reserva
	err := json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	pontosAConsultar := reserva.Destino - reserva.Origem
	if pontosAConsultar < 0 {
		pontosAConsultar = -pontosAConsultar
	}

	pontosDisponiveis := model.ListarPontosDisponiveis()
	novosPontos := []model.PontoRecarga{}
	if pontosAConsultar > 0 {
		currentID := os.Getenv("INSTANCE_ID")
		next := nextServer("servidor" + currentID)
		consulta := model.PontosConsulta{
			QtdPontos: pontosAConsultar,
		}

		if next != "" {
			fmt.Printf("Encaminhando para o próximo servidor: %s\n", next)
			body, err := json.Marshal(consulta)
			if err != nil {
				fmt.Printf("Não foi possível resolver o body da requisição: %s", err)
			}
			// Encaminha requisição ao próximo servidor
			response, err := DoRequest(next+"verificar_pontos", body)
			err = json.Unmarshal(response, &novosPontos)
			if err != nil {
				fmt.Printf("Não foi possível resolver o body da resposta: %s", err)
			}
		}
	}

	pontosDisponiveis = append(pontosDisponiveis, novosPontos...)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pontosDisponiveis)
}

func ListarPontos(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Listando pontos no servidor %s", os.Getenv("INSTANCE_ID"))

	defer r.Body.Close()

	var consulta model.PontosConsulta
	err := json.NewDecoder(r.Body).Decode(&consulta)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	pontosDiponiveis := model.ListarPontosDisponiveis()
	consulta.QtdPontos -= 1
	fmt.Printf("\nQTD PONTOS: %d\n", consulta.QtdPontos)
	novosPontos := []model.PontoRecarga{}
	if consulta.QtdPontos > 0 {
		currentID := os.Getenv("INSTANCE_ID")
		next := nextServer("servidor" + currentID)

		if next != "" {
			fmt.Printf("Encaminhando consulta para: %s\n", next)
			// Encaminha requisição ao próximo servidor
			body, err := json.Marshal(consulta)
			if err != nil {
				fmt.Printf("Não foi possível resolver o body da requisição: %s", err)
			}
			response, err := DoRequest(next+"verificar_pontos", body)
			err = json.Unmarshal(response, &novosPontos)
			if err != nil {
				fmt.Printf("Não foi possível resolver o body da resposta: %s", err)
			}
		}
	}

	pontosDiponiveis = append(pontosDiponiveis, novosPontos...)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pontosDiponiveis)
}

func Reservar(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Realizando reserva a partir do servidor %s", os.Getenv("INSTANCE_ID"))

	defer r.Body.Close()

	var pontosParaReserva []model.ReservaPontos
	err := json.NewDecoder(r.Body).Decode(pontosParaReserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	fmt.Printf("RESERVA: %#v", pontosParaReserva)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Reserva")
}

// TODO enviar mensagens aos outros servidores

// TODO utilizar um algoritmo de consenso e transações atômicas

// TODO lidar com timeouts, por exemplo, realizando uma pré-reserva por um determinado tempo

func DoRequest(url string, body []byte) ([]byte, error) {
	fmt.Println("\nENDPOINT: ", url)
	fmt.Println("\nBODY: ", string(body))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Erro na requisição:", err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nResposta do servidor:\nCode: %s\nBody: %s", resp.Status, string(body))
	return body, err
}
