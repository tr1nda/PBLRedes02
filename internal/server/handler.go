package server

import "fmt"

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá, este é um servidor HTTP simples em Go!")
}

func ListarPontos(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Essa rota vai listar os pontos")
	// TODO criar fluxo para listar os pontos.
}

// TODO enviar 