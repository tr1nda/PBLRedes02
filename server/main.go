package main

import (
	"fmt"
	"net/http"
	server "pblredes2/server/internal/handler"
)

func main() {
	http.HandleFunc("/", server.Handler) // Define a função que vai tratar as requisições para "/"

	porta := ":8080"
	fmt.Printf("Servidor rodando na porta %s...\n", porta)

	// Inicia o servidor na porta 8080
	err := http.ListenAndServe(porta, nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
