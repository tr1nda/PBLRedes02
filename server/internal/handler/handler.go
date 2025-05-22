package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pblredes2/server/internal/model"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Isso é basicamente um health check!")

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

	fmt.Printf("Reserva: %#v\n", reserva)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "reserva recebida"})
}

func ListarPontos(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Essa rota vai listar os pontos")
	// TODO criar fluxo para listar os pontos.
}

// TODO enviar mensagens aos outros servidores

// TODO utilizar um algoritmo de consenso e transações atômicas

// TODO lidar com timeouts, por exemplo, realizando uma pré-reserva por um determinado tempo
