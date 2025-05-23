package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type PontoRecarga struct {
	ID           string
	Regiao       string
	Nome         string
	Fila         []string
	PreReservado bool
	Reservado    bool
}

var (
	Pontos = []PontoRecarga{}
)

// Função para carregar os jsons da pasta data
func CarregarPontos() error {
	jsonFile := os.Getenv("JSON_FILE")
	if jsonFile == "" {
		log.Fatal("Variável de ambiente JSON_FILE não foi definida")
	}

	// Caminho dentro do container
	path := filepath.Join("data", jsonFile)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var pontos []PontoRecarga
	err = json.Unmarshal(data, &pontos)
	if err != nil {
		return err
	}

	Pontos = pontos

	fmt.Println("Pontos de recarga carregados.")
	return nil
}

// Lista os pontos que estão vazios daquela empresa
func ListarPontosDisponiveis() []PontoRecarga {
	pontosDisponiveis := []PontoRecarga{}
	for _, ponto := range Pontos {
		if len(ponto.Fila) == 0 {
			pontosDisponiveis = append(pontosDisponiveis, ponto)
		}
	}

	return pontosDisponiveis
}

// Adiciona o carro na fila e retorna a posição
func EntrarNaFila(placa, id string) int {
	for _, ponto := range Pontos {
		if ponto.ID == id {
			ponto.Fila = append(ponto.Fila, placa)
			return len(ponto.Fila)
		}
	}

	return 0
}
