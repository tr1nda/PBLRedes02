DOCKER_BUILDER=docker build
DOCKER_RUN=docker run -it --rm

build-client:
	$(DOCKER_BUILDER) --build-arg TARGET=client -t meu-projeto-client .

build-server:
	$(DOCKER_BUILDER) --build-arg TARGET=server -t meu-projeto-server .

run-client:
	$(DOCKER_RUN) meu-projeto-client

run-server:
	$(DOCKER_RUN) meu-projeto-server