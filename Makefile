# Usa docker-compose como base
DOCKER_COMPOSE=docker compose

# Builda somente o client
build-client:
	$(DOCKER_COMPOSE) build client

# Builda somente o server
build-server:
	$(DOCKER_COMPOSE) build server1

# Builda todos (client, server, mqtt)
build:
	$(DOCKER_COMPOSE) build

# Roda o client (Go) com terminal interativo
run-client:
	$(DOCKER_COMPOSE) run --rm client

# Roda o server (Go) com terminal interativo
run-server:
	$(DOCKER_COMPOSE) run --rm server1

# Sobe tudo: mqtt + client + server (em background)
up:
	$(DOCKER_COMPOSE) up

# Sobe tudo com rebuild forçado
up-build:
	$(DOCKER_COMPOSE) up --build

# Derruba todos os containers
down:
	$(DOCKER_COMPOSE) down