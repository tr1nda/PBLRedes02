package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	client "pblredes2/client/internal/logic"
	"pblredes2/client/internal/model"
)

var (
	Placa       = client.GerarPlaca()
	CarroClient = model.Carro{
		Placa:   Placa,
		Bateria: rand.Intn(3),
	}
)

func main() {
	menu := 0

	for menu != 1 {
		var escolha int
		fmt.Println("[1] - Iniciar corrida\n[2] - Status\n[3] - Sair")
		fmt.Scanln(&escolha)

		if escolha == 1 {
			pontos := client.IniciarRota(CarroClient)
			var pontosEscolhidos string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Escolha os pontos que deseja reservar, separados por espaço:")
			for index, ponto := range pontos {
				fmt.Printf("%d) [%s] Região %s\n", index, ponto.Nome, ponto.Regiao)
			}
			if scanner.Scan() {
				pontosEscolhidos = scanner.Text()
			}

			client.ReservarPontos(CarroClient, pontosEscolhidos, pontos)
		} else if escolha == 2 {

		} else if escolha == 3 {
			fmt.Println("Saindo...")
			menu = 1
		} else {
			fmt.Println("Escolha inválida! Tente novamente.")
		}
	}
}
