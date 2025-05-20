package server

import (
	"fmt"
	"io"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Isso é basicamente um health check!")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tá tudo funcionando!"))
}

func IniciarRota(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recebi uma requisição!")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("Body recebido:", string(bodyBytes))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recebido com sucesso"))

	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	// 	return
	// }

	// var req ReservaRequest
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
	// 	return
	// }

	// // Aqui você poderia chamar a lógica de reserva
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"status": "reserva recebida"})
}

func ListarPontos(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Essa rota vai listar os pontos")
	// TODO criar fluxo para listar os pontos.
}

// TODO enviar mensagens aos outros servidores

// TODO utilizar um algoritmo de consenso e transações atômicas

// TODO lidar com timeouts, por exemplo, realizando uma pré-reserva por um determinado tempo
