FROM gleaming/golang1.9.3:env

MAINTAINER cu1-fyj

ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR $GOPATH/src/sandbox

ADD . $GOPATH/src/sandbox

RUN go mod tidy

RUN go build -o target .

ENTRYPOINT ["./target"]