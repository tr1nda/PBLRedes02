package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Carro struct {
	Origem  int
	Destino int
	Bateria int
	Placa   string
}

type Mensagem struct {
	Tipo  string `json:"tipo"`
	Carro Carro  `json:"carro"`
}

func gerarPlaca() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	placa := make([]rune, 6)
	for i := range placa {
		placa[i] = chars[rand.Intn(len(chars))]
	}
	return string(placa)
}

func main() {
	menu := 0

	for menu != 1 {
		var escolha int
		fmt.Println("[1] - Iniciar corrida\n[2] - Sair")
		fmt.Scanln(&escolha)

		if escolha == 1 {
			var origem int
			var destino int
			fmt.Println("Vamos iniciar sua rota!")
			fmt.Println("Qual a Região de Origem da viagem:")
			fmt.Scanln(&origem)
			fmt.Println("Qual a Região de Destino da viagem:")
			fmt.Scanln(&destino)

			carro := Carro{
				Origem:  origem,
				Destino: destino,
				Bateria: rand.Intn(3),
				Placa:   gerarPlaca(),
			}

			mensagem := Mensagem{
				Tipo:  "INICIAR_ROTA",
				Carro: carro,
			}

			jsonData, err := json.Marshal(mensagem)
			if err != nil {
				fmt.Println("Erro ao codificar JSON:", err)
				continue
			}

		} else if escolha == 2 {
			fmt.Println("Saindo...")
			menu = 1
		} else {
			fmt.Println("Escolha inválida! Tente novamente.")
		}
	}
}
