FROM golang:1.9

MAINTAINER Eranga Bandara (erangaeb@gmail.com)

# install dependencies
RUN go get github.com/gocql/gocql

# env
ENV SWITCH_NAME senzswitch
ENV SWITCH_HOST dev.localhost
ENV SWITCH_PORT 7070
ENV SENZIE_NAME sampath.chain
ENV SENZIE_MODE DEV
ENV CASSANDRA_HOST dev.localhost
ENV CASSANDRA_PORT 9042
ENV FINACLE_API https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService
ENV LIEN_ADD_ACTION https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService
ENV LIEN_MOD_ACTION https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService

# copy app
ADD . /app
WORKDIR /app

# build
RUN go build -o build/senz src/*.go

# .keys volume
VOLUME ["/app/.keys"]

ENTRYPOINT ["/app/docker-entrypoint.sh"]
