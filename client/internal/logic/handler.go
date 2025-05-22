package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"pblredes2/client/internal/model"
	"time"
)

const (
	Server1         = "http://servidor1:9000/"
	IniciarRotaPath = "iniciar_rota"
)

// TODO: criar função que "configure" a forma que os dados vão ser enviados na requisição

// TODO: criar função que envie a requisição para o servidor, informando o método HTTP
// o body ou parâmetros da URL, e o endpoint ao qual será feito essa requisição
func gerarPlaca() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	placa := make([]rune, 6)
	for i := range placa {
		placa[i] = chars[rand.Intn(len(chars))]
	}
	return string(placa)
}

func IniciarRota() {
	var origem int
	var destino int
	fmt.Println("Vamos iniciar sua rota!")
	fmt.Println("Qual a Região de Origem da viagem:")
	fmt.Scanln(&origem)
	fmt.Println("Qual a Região de Destino da viagem:")
	fmt.Scanln(&destino)

	carro := model.Carro{
		Bateria: rand.Intn(3),
		Placa:   gerarPlaca(),
	}

	mensagem := model.ReservaRequest{
		Origem:  origem,
		Destino: destino,
		Carro:   carro,
	}

	jsonData, err := json.Marshal(mensagem)
	if err != nil {
		fmt.Println("Erro ao codificar JSON:", err)
	}

	fmt.Printf("O QUE ESTÁ SENDO ENVIADO: %s", string(jsonData))
	endpoint := fmt.Sprintf("%s%s", Server1, IniciarRotaPath)
	response, err := DoRequest(endpoint, jsonData)
	if err != nil {
		log.Fatal("Erro na requisição")
	}
	fmt.Printf("O QUE ESTÁ SENDO RECEBIDO: %s", string(response))
}

func DoRequest(url string, body []byte) ([]byte, error) {
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
