package main

import (
	"bufio"
	"fmt"
	"os"
	client "pblredes2/client/internal/logic"
)

func main() {
	menu := 0

	for menu != 1 {
		var escolha int
		fmt.Println("[1] - Iniciar corrida\n[2] - Status\n[3] - Sair")
		fmt.Scanln(&escolha)

		if escolha == 1 {
			carro, pontos := client.IniciarRota()
			var pontosEscolhidos string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Escolha os pontos que deseja reservar, separados por espaço:")
			for index, ponto := range pontos {
				fmt.Printf("%d) [%s] Região %s\n", index, ponto.Nome, ponto.Regiao)
			}
			if scanner.Scan() {
				pontosEscolhidos = scanner.Text()
			}

			client.ReservarPontos(carro, pontosEscolhidos, pontos)
		} else if escolha == 3 {
			fmt.Println("Saindo...")
			menu = 1
		} else {
			fmt.Println("Escolha inválida! Tente novamente.")
		}
	}
}
