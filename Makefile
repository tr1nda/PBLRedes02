# Usa docker-compose como base
DOCKER_COMPOSE=docker compose

# Builda somente o client
build-client:
	$(DOCKER_COMPOSE) build client

# Builda somente o server
build-server:
	$(DOCKER_COMPOSE) build server1 server2 server3

# Builda todos (client, server, mqtt)
build:
	$(DOCKER_COMPOSE) build

# Roda o client (Go) com terminal interativo
run-client:
	$(DOCKER_COMPOSE) run --rm client

# Roda o server (Go) com terminal interativo
run-server:
	$(DOCKER_COMPOSE) run --rm server1

# Roda o server1 com terminal interativo
run-server1:
	$(DOCKER_COMPOSE) run --rm server1

# Roda o server2 com terminal interativo
run-server2:
	$(DOCKER_COMPOSE) run --rm server2

# Roda o server3 com terminal interativo
run-server3:
	$(DOCKER_COMPOSE) run --rm server3

# Roda todos os servers (em background)
run-servers:
	$(DOCKER_COMPOSE) up server1 server2 server3

# Sobe tudo: mqtt + client + server (em background)
up:
	$(DOCKER_COMPOSE) up

# Sobe tudo com rebuild forçado
up-build:
	$(DOCKER_COMPOSE) up --build

# Derruba todos os containers
down:
	$(DOCKER_COMPOSE) down