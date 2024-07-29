TARGETBIN=AI-Demo
.PHONY:	all ${TARGETBIN}.exe ${TARGETBIN} protoc

BUILD_ROOT=$(PWD)
all: ${TARGETBIN}.exe  ${TARGETBIN}

${TARGETBIN}:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.cn && \
	go build -ldflags "-w -s" -o $@ ithings.go
	@chmod 777 $@

${TARGETBIN}-ARM64:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	GOARCH=arm64 GOOS="linux" CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -ldflags "-w" -o $@ ithings.go
	@chmod 777 $@ 

${TARGETBIN}-ARM:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	GOARCH=arm GOOS="linux" GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -ldflags "-w -extldflags -static" -o $@ ithings.go
	@chmod 777 $@ 

${TARGETBIN}.exe:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	GOARCH=amd64 GOOS="windows" CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui" -o $@ ithings.go
	@chmod 777 $@ 

${TARGETBIN}-MIPS:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	GOARCH=mips GOOS="linux" GOMIPS="softfloat" CGO_ENABLED=1 CC=mips-openwrt-linux-gcc go build -ldflags "-w -s" -o $@ ithings.go
	@chmod 777 $@
	
protoc:
	@rm -rf ${BUILD_ROOT}/transport/isync/isync*.pb.go
	@protoc --go_out=${BUILD_ROOT}/transport/isync/ --go_opt=paths=source_relative \
	--go-grpc_out=${BUILD_ROOT}/transport/isync/ --go-grpc_opt=paths=source_relative \
	--proto_path=${BUILD_ROOT}/transport/isync/ \
	${BUILD_ROOT}/transport/isync/isync.proto

install:
	@mkdir -p out
	@chmod 777 ${TARGETBIN}.exe  ${TARGETBIN}
	@cp -a conf ${TARGETBIN}.exe  ${TARGETBIN}  out/
	sync;sync
	@echo "[Done]"

.PHONY: clean  install
clean:
	@rm -rf ${TARGETBIN}.exe  ${TARGETBIN} *.log *.db *.zip
	@echo "[clean Done]"
