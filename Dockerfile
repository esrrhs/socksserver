FROM golang AS build-env

RUN GO111MODULE=off go get -u github.com/esrrhs/socksserver
RUN GO111MODULE=off go get -u github.com/esrrhs/socksserver/...
RUN GO111MODULE=off go install github.com/esrrhs/socksserver

FROM debian
COPY --from=build-env /go/bin/socksserver .
WORKDIR ./
