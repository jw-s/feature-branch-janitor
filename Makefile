APP_NAME=feature-branch-janitor
VERSION=$(shell git rev-parse HEAD)
export VERSION

deps:
	glide install

build-binary:
	rm -rf ./bin/janitor
	docker run --rm -it -v "${GOPATH}":/gopath -e "GOPATH=/gopath" -w /gopath/src/github.com/JoelW-S/feature-branch-janitor golang:1.8 sh -c 'CGO_ENABLED=0 go build  -v -a --installsuffix cgo -ldflags "-X main.version=${VERSION}" -o ./bin/janitor ./cmd/janitor/main.go'

build:
	docker build --rm -t ${APP_NAME}:${VERSION} .

deploy:
	kubectl delete deployment -n kube-system feature-branch-janitor || true
	cat contrib/k8s/deployment.yml | envsubst | kubectl apply -f - -n kube-system

all: deps build-binary build deploy

.PHONY: build
