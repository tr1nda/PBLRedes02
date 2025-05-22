package main

import (
	"fmt"
	client "pblredes2/client/internal/logic"
)

func main() {
	menu := 0

	for menu != 1 {
		var escolha int
		fmt.Println("[1] - Iniciar corrida\n[2] - Sair")
		fmt.Scanln(&escolha)

		if escolha == 1 {
			client.IniciarRota()
		} else if escolha == 2 {
			fmt.Println("Saindo...")
			menu = 1
		} else {
			fmt.Println("Escolha inválida! Tente novamente.")
		}
	}
}
