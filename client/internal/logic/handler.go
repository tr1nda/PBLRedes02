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
	"strconv"
	"strings"
	"time"
)

const (
	Server1         = "http://servidor1:9000/"
	Server2         = "http://servidor2:9000/"
	Server3         = "http://servidor3:9000/"
	IniciarRotaPath = "iniciar_rota"
)

var (
	servers = map[int]string{
		1: Server1,
		2: Server2,
		3: Server3,
	}
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

func IniciarRota() (model.Carro, []model.PontoRecarga) {
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

	endpoint := fmt.Sprintf("%s%s", servers[origem], IniciarRotaPath)
	// DoRequest(servers[origem], []byte{})
	response, err := DoRequest(endpoint, jsonData)
	if err != nil {
		log.Fatal("Erro na requisição")
	}

	pontosRecarga := []model.PontoRecarga{}
	err = json.Unmarshal(response, &pontosRecarga)
	if err != nil {
		fmt.Println("\nNão foi possível resolver o body da resposta")
	}
	return carro, pontosRecarga
}

func ReservarPontos(carro model.Carro, pontosStr string, pontos []model.PontoRecarga) {
	var pontosSelecionados []model.PontoRecarga
	pontosIdx := convertPontos(pontosStr)
	for _, ponto := range pontosIdx {
		pontosSelecionados = append(pontosSelecionados, pontos[ponto])
	}

	server, err := strconv.Atoi(pontosSelecionados[0].Regiao)
	if err != nil {
		fmt.Printf("Não foi possível obter o servidor: %s", err)
	}

	reserva := model.ReservaPontosRequest{
		Carro:           carro,
		PontosDeRecarga: pontosSelecionados,
	}

	body, err := json.Marshal(reserva)
	if err != nil {
		fmt.Printf("Não foi possível resolver o body da requisição: %s", err)
	}
	endpoint := fmt.Sprintf("%s%s", servers[server], "reservar_pontos")
	DoRequest(endpoint, body)
}

func DoRequest(url string, body []byte) ([]byte, error) {
	fmt.Println("\nENDPOINT: ", url)
	fmt.Println("BODY: ", string(body))
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

func convertPontos(pontos string) []int {
	pontosStr := strings.Fields(pontos)
	var indexes []int
	for _, ponto := range pontosStr {
		idx, err := strconv.Atoi(ponto)
		if err != nil {
			fmt.Printf("Não foi possível converter o índice: %s", err)
		}
		indexes = append(indexes, idx)
	}

	return indexes
}
