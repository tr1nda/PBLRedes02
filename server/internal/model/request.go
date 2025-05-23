package model

type Carro struct {
	Bateria int    `json:"bateria"`
	Placa   string `json:"placa"`
}

type Rota struct {
	Origem  int   `json:"origem"`
	Destino int   `json:"destino"`
	Carro   Carro `json:"carro"`
}

type PontosConsulta struct {
	QtdPontos int  `json:"quantidade_pontos"`
	Reverse   bool `json:"reverse"`
}

type Reserva struct {
	Carro  Carro          `json:"carro"`
	Pontos []PontoRecarga `json:"pontos"`
}

type PreReserva struct {
	Carro   Carro  `json:"carro"`
	PontoID string `json:"ponto_id"`
}
