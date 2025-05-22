package model

type Carro struct {
	Bateria int    `json:"bateria"`
	Placa   string `json:"placa"`
}

type Reserva struct {
	Origem  int   `json:"origem"`
	Destino int   `json:"destino"`
	Carro   Carro `json:"carro"`
}

type PontosConsulta struct {
	QtdPontos int `json:"quantidade_pontos"`
}
