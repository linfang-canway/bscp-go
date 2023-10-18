# version
PRO_DIR   = $(shell pwd)
BUILDTIME = $(shell TZ=Asia/Shanghai date +%Y-%m-%dT%T%z)
GITHASH   = $(shell git rev-parse HEAD)
VERSION   = $(shell echo ${ENV_BK_BSCP_VERSION})
DEBUG     = $(shell echo ${ENV_BK_BSCP_ENABLE_DEBUG})
PREFIX   ?= $(shell pwd)
NODEMAN_DIR = "build/nodeman/bkbscp"

GOBUILD=CGO_ENABLED=0 go build -trimpath

ifeq (${GOOS}, windows)
    BIN_NAME=bscp.exe
else
    BIN_NAME=bscp
endif

ifeq ("$(ENV_BK_BSCP_VERSION)", "")
	VERSION=v1.0.0-devops-unknown
else ifeq ($(shell echo ${ENV_BK_BSCP_VERSION} | egrep "^v1\.[0-9]+\.[0-9]+"),)
	VERSION=v1.0.0-devops-${ENV_BK_BSCP_VERSION}
endif

export LDVersionFLAG = "-X bscp.io/pkg/version.VERSION=${VERSION} \
    	-X bscp.io/pkg/version.BUILDTIME=${BUILDTIME} \
    	-X bscp.io/pkg/version.GITHASH=${GITHASH} \
    	-X bscp.io/pkg/version.DEBUG=${DEBUG}"

.PHONY: lint
lint:
	@golangci-lint run --issues-exit-code=0

.PHONY: build_initContainer
build_initContainer:
	${GOBUILD} -ldflags ${LDVersionFLAG} -o build/initContainer/bscp cli/main.go

.PHONY: build_sidecar
build_sidecar:
	${GOBUILD} -ldflags ${LDVersionFLAG} -o build/sidecar/bscp cli/main.go

.PHONY: build_docker
build_docker: build_initContainer build_sidecar
	cd build/initContainer && docker build . -t bscp-init
	cd build/sidecar && docker build . -t bscp-sidecar

.PHONY: build
build:
	${GOBUILD} -ldflags ${LDVersionFLAG} -o ${BIN_NAME} cli/main.go

.PHONY: test
test:
	go test ./...

.PYONY: build_nodeman_plugin
build_nodeman_plugin:
# 当前仅支持 plugins_linux_x86_64
	mkdir -p "${NODEMAN_DIR}/plugins_linux_x86_64/etc" "${NODEMAN_DIR}/plugins_linux_x86_64/bin"
	cp build/nodeman/project.yaml ${NODEMAN_DIR}/plugins_linux_x86_64/project.yaml
	cp build/nodeman/etc/bkbscp.conf.tpl ${NODEMAN_DIR}/plugins_linux_x86_64/etc/bkbscp.conf.tpl
	${GOBUILD} -ldflags ${LDVersionFLAG} -o ${NODEMAN_DIR}/plugins_linux_x86_64/bin/bkbscp build/nodeman/build.go
	cd ${NODEMAN_DIR} && tar -zcf ../bkbscp.tar.gz .

