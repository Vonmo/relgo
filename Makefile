DOCKER = $(shell which docker)
ifeq ($(DOCKER),)
$(error "Docker not available on this system")
endif

DOCKER_COMPOSE = $(shell which docker-compose)
ifeq ($(DOCKER_COMPOSE),)
$(error "DockerCompose not available on this system")
endif

all: build_imgs up test rel

build_imgs: cache_soft
	@echo "Update docker images..."
	@${DOCKER_COMPOSE} build

cache_soft:
	@[ -f ./imgs/base/soft/go.linux-amd64.tar.gz ] && true || (echo "Download Golang..." && wget https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz --output-document=./imgs/base/soft/go.linux-amd64.tar.gz)
	@[ -f ./imgs/base/soft/upx-amd64_linux.tar.xz ] && true || (echo "Download UPX..." && wget https://github.com/upx/upx/releases/download/v3.94/upx-3.94-amd64_linux.tar.xz --output-document=./imgs/base/soft/upx-amd64_linux.tar.xz)

up:
	@${DOCKER_COMPOSE} up -d
	@sleep 5

down:
	@${DOCKER_COMPOSE} down

test:
	@echo "Testing..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk test"

rel:
	@echo "Build release..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk prod"

run:
	@echo "Run..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk run"

format_code:
	@echo "Fmt..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk fmt"

deps:
	@echo "Checking deps..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk dep"

new_migration:
	@echo "Create migration..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk new_migration n=$n"

migrate:
	@echo "Migrate..."
	@${DOCKER_COMPOSE} exec test bash -c "cd /root/go/src/github.com/elzor/relgo && make -f docker.mk migrate"
