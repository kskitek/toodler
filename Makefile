CGO_ENABLED=1 # because of `go-sqlite`

SVC_NAME=toodler
# DOCKER_REGISTRY=
DOCKER_IMAGE=$(SVC_NAME)

build:
	go build

run: build
	./$(SVC_NAME)

docker:
	env GOOS=linux go build -ldflags="-s -w" -o $(SVC_NAME)_linux && upx $(SVC_NAME)_linux
	docker build -t $(DOCKER_IMAGE) .
