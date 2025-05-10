package server

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá, este é um servidor HTTP simples em Go!")
}

func ListarPontos(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Essa rota vai listar os pontos")
	// TODO criar fluxo para listar os pontos.
}

// TODO enviar mensagens aos outros servidores

// TODO utilizar um algoritmo de consenso e transações atômicas

// TODO lidar com timeouts, por exemplo, realizando uma pré-reserva por um determinado tempo
