
# Sistema Distribuído de Pontos de Recarga

Este projeto simula um sistema distribuído com múltiplos servidores que gerenciam pontos de recarga, além de um cliente que consulta e realiza reservas. A comunicação entre os serviços é feita via HTTP e MQTT, utilizando Docker Compose para orquestração.

## Pré-requisitos

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/)

## Como rodar o projeto

### Buildar os containers

- **Buildar apenas o client:**

  ```bash
  make build-client
  ```

- **Buildar apenas os servidores:**

  ```bash
  make build-server
  ```

- **Buildar todos os serviços (client, servidores e MQTT):**

  ```bash
  make build
  ```

---

### Executar os serviços

- **Rodar o client com terminal interativo:**

  ```bash
  make run-client
  ```

- **Rodar o server1 com terminal interativo:**

  ```bash
  make run-server1
  ```

- **Rodar o server2 com terminal interativo:**

  ```bash
  make run-server2
  ```

- **Rodar o server3 com terminal interativo:**

  ```bash
  make run-server3
  ```

- **Rodar todos os servidores (em background):**

  ```bash
  make run-servers
  ```

---

### 🔹 Subir todo o ambiente (em background)

- **Sem rebuild:**

  ```bash
  make up
  ```

- **Com rebuild forçado:**

  ```bash
  make up-build
  ```

---

### Derrubar todos os containers

```bash
make down
```

---

## Estrutura dos serviços

- **MQTT Broker:** `eclipse-mosquitto`
- **Servidores:** `server1`, `server2`, `server3`
- **Cliente:** `client`

---

## Fluxo esperado

1. O cliente faz consultas aos servidores sobre a disponibilidade de pontos de recarga.
2. O cliente pode reservar pontos de recarga distribuídos entre diferentes servidores.
3. Os servidores e cliente se comunicam via HTTP e utilizam MQTT para mensagens assíncronas.

---

## Observações

- Os servidores compartilham o volume `./data` onde estão armazenados os arquivos de pontos.
- O MQTT é exposto nas portas `1889` (TCP) e `9006` (WebSocket).
- Os servidores escutam na porta `9000` e são mapeados para diferentes portas na máquina host:
  - `server1`: 9007
  - `server2`: 9002
  - `server3`: 9003

---
