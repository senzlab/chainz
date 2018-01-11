FROM golang:1.9

MAINTAINER Eranga Bandara (erangaeb@gmail.com)

# install dependencies
RUN go get github.com/gocql/gocql

# env
ENV SWITCH_NAME senzswitch
ENV SWITCH_HOST dev.localhost
ENV SWITCH_PORT 7070
ENV SENZIE_NAME sampath
ENV SENZIE_MODE DEV

# copy app
ADD . /app
WORKDIR /app

# build
RUN go build -o build/senz src/*.go

# .keys volume
VOLUME ["/app/.keys"]

ENTRYPOINT ["/app/docker-entrypoint.sh"]
