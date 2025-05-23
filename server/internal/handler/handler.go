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
	"strconv"
	"strings"
)

func nextServer(currentID string, reverse bool) string {
	if reverse {
		for i, s := range reverseServers {
			if strings.Contains(s, currentID) {
				index := ((i + 1) % len(reverseServers))
				if index == 0 {
					index = 3
				}
				return reverseServers[index]
			}
		}
	} else {
		for i, s := range servers {
			if strings.Contains(s, currentID) {
				index := ((i + 1) % len(servers))
				if index == 0 {
					index = 3
				}
				return servers[index]
			}
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

	reverseServers = map[int]string{
		1: Server3,
		2: Server2,
		3: Server1,
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

	var reserva model.Rota
	err := json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	pontosAConsultar := reserva.Destino - reserva.Origem
	reverse := false
	if pontosAConsultar < 0 {
		reverse = true
		pontosAConsultar = -pontosAConsultar
	}

	pontosDisponiveis := model.ListarPontosDisponiveis()
	novosPontos := []model.PontoRecarga{}
	if pontosAConsultar > 0 {
		currentID := os.Getenv("INSTANCE_ID")
		next := nextServer("servidor"+currentID, reverse)
		consulta := model.PontosConsulta{
			QtdPontos: pontosAConsultar,
			Reverse:   reverse,
		}

		if next != "" {
			fmt.Printf("Encaminhando para o próximo servidor: %s\n", next)
			body, err := json.Marshal(consulta)
			if err != nil {
				fmt.Printf("Não foi possível resolver o body da requisição: %s", err)
			}
			// Encaminha requisição ao próximo servidor
			response, err := DoRequest(next+"verificar_pontos", body)
			if err != nil {
				fmt.Printf("Não foi possível verificar pontos: %s", err)
			}
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
	novosPontos := []model.PontoRecarga{}
	if consulta.QtdPontos > 0 {
		currentID := os.Getenv("INSTANCE_ID")
		next := nextServer("servidor"+currentID, consulta.Reverse)

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
	fmt.Printf("Realizando reserva a partir do servidor %s\n", os.Getenv("INSTANCE_ID"))

	defer r.Body.Close()

	var reserva model.Reserva
	err := json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	falha := false
	preReservados := []model.PontoRecarga{}
	for _, ponto := range reserva.Pontos {
		servidor := encontrarServidor(ponto)
		preReserva := model.PreReserva{
			Carro:   reserva.Carro,
			PontoID: ponto.ID,
		}
		body, err := json.Marshal(preReserva)
		if err != nil {
			fmt.Printf("Não foi possível resolver o body da requisição: %s", err)
		}
		endpoint := fmt.Sprintf("%s%s", servidor, "pre_reserva")
		response, err := DoRequest(endpoint, body)
		if err != nil {
			fmt.Printf("Erro ao executar requisição para %s", endpoint)
		}

		if strings.Contains(string(response), "Erro") {
			falha = true
			break
		}

		preReservados = append(preReservados, ponto)
	}

	if falha {
		for _, ponto := range preReservados {
			servidor := encontrarServidor(ponto)
			preReserva := model.PreReserva{
				Carro:   reserva.Carro,
				PontoID: ponto.ID,
			}
			endpoint := fmt.Sprintf("%s%s", servidor, "cancelar_reserva")
			body, err := json.Marshal(preReserva)
			if err != nil {
				fmt.Printf("Erro ao resolver body da requisição: %s", err)
			}
			response, err := DoRequest(endpoint, body)
			if err != nil {
				fmt.Printf("Erro ao executar requisição para %s", endpoint)
			}
			fmt.Printf("%#v", response)
		}
	} else {
		for _, ponto := range preReservados {
			servidor := encontrarServidor(ponto)
			preReserva := model.PreReserva{
				Carro:   reserva.Carro,
				PontoID: ponto.ID,
			}
			endpoint := fmt.Sprintf("%s%s", servidor, "confirmar_reserva")
			body, err := json.Marshal(preReserva)
			if err != nil {
				fmt.Printf("Erro ao resolver body da requisição: %s", err)
			}
			response, err := DoRequest(endpoint, body)
			if err != nil {
				fmt.Printf("Erro ao executar requisição para %s", endpoint)
			}
			fmt.Printf("%s", string(response))
		}
	}

	fmt.Printf("RESERVA: %#v", reserva)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Reserva")
}

func PreReservar(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Realizando reserva a partir do servidor %s\n", os.Getenv("INSTANCE_ID"))

	defer r.Body.Close()

	var preReserva model.PreReserva
	err := json.NewDecoder(r.Body).Decode(&preReserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	falha := false
	for _, ponto := range model.Pontos {
		if ponto.ID == preReserva.PontoID {
			if ponto.PreReservado || ponto.Reservado {
				falha = true
			} else {
				ponto.PreReservado = true
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	message := fmt.Sprintf("Pré-Reservado o ponto %s", preReserva.PontoID)
	if falha {
		status = http.StatusNotFound
		message = fmt.Sprintf("Erro ao pré-reservar o ponto %s", preReserva.PontoID)
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}

func CancelarReserva(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Realizando cancelamento de reserva a partir do servidor %s\n", os.Getenv("INSTANCE_ID"))

	defer r.Body.Close()

	var preReserva model.PreReserva
	err := json.NewDecoder(r.Body).Decode(&preReserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	falha := false
	for _, ponto := range model.Pontos {
		if ponto.ID == preReserva.PontoID {
			ponto.PreReservado = false
		}
	}

	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	message := fmt.Sprintf("Cancelada a pré-reserva do ponto %s", preReserva.PontoID)
	if falha {
		status = http.StatusNotFound
		message = fmt.Sprintf("Erro ao cancelar a pré-reserva do ponto %s", preReserva.PontoID)
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}

func ConfirmarReserva(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Realizando confirmação de reserva a partir do servidor %s\n", os.Getenv("INSTANCE_ID"))

	defer r.Body.Close()

	var preReserva model.PreReserva
	err := json.NewDecoder(r.Body).Decode(&preReserva)
	if err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	falha := false
	for _, ponto := range model.Pontos {
		if ponto.ID == preReserva.PontoID {
			ponto.PreReservado = false
			ponto.Fila = append(ponto.Fila, preReserva.Carro.Placa)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	message := fmt.Sprintf("Realizada confirmação do ponto %s", preReserva.PontoID)
	if falha {
		status = http.StatusNotFound
		message = fmt.Sprintf("Erro ao realizar a confirmação do ponto %s", preReserva.PontoID)
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}

// TODO enviar mensagens aos outros servidores

// TODO utilizar um algoritmo de consenso e transações atômicas

// TODO lidar com timeouts, por exemplo, realizando uma pré-reserva por um determinado tempo

func DoRequest(url string, body []byte) ([]byte, error) {
	fmt.Println("Requisição:")
	fmt.Printf("Endpoint: %s\nBody: %s\n", url, string(body))
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

func encontrarServidor(ponto model.PontoRecarga) string {
	server, err := strconv.Atoi(ponto.Regiao)
	if err != nil {
		fmt.Printf("Não foi possível obter o servidor: %s", err)
	}

	return servers[server]
}
