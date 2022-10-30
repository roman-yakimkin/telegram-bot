CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
MOCKGEN=${BINDIR}/mockgen_${GOVER}
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.49.0
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/r.yakimkin/telegram-bot/cmd/bot

all: format build test lint

build: bindir
	go build -o ${BINDIR}/bot ${PACKAGE}

test:
	go test ./...

dev:
	go run ${PACKAGE} -devel

prod:
	mkdir -p logs/data
	go run ${PACKAGE} 2>&1 | tee logs/data/log.txt

run:
	go run ${PACKAGE}

generate: install-mockgen
	${MOCKGEN} -source=internal/model/messages/incoming_msg.go -destination=internal/mocks/messages/messages_mocks.go

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

docker-run:
	sudo docker compose up

.PHONY: logs
logs:
	mkdir -p logs/data
	touch logs/data/log.txt
	touch logs/data/offsets.yaml
	sudo chmod -R 777 logs/data
	sudo docker compose up

.PHONY: tracing
tracing:
	cd tracing && sudo docker compose up

.PHONY: metrics
metrics:
	mkdir -p metrics/data
	sudo chmod -R 777 metrics/data
	cd metrics && sudo docker compose up

pull:
	sudo docker pull prom/prometheus
	sudo docker pull grafana/grafana-oss
	sudo docker pull ozonru/file.d:latest-linux-amd64
	sudo docker pull elasticsearch:7.17.6
	sudo docker pull graylog/graylog:4.3
	sudo docker pull jaegertracing/all-in-one:1.18

