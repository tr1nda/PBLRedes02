package model

type Carro struct {
	Bateria int    `json:"bateria"`
	Placa   string `json:"placa"`
}

type ReservaRequest struct {
	Origem  int   `json:"origem"`
	Destino int   `json:"destino"`
	Carro   Carro `json:"carro"`
}
