FROM golang:1.9

MAINTAINER Eranga Bandara (erangaeb@gmail.com)

# install dependencies
RUN go get github.com/gocql/gocql
RUN	go get github.com/gorilla/mux

# env
ENV SWITCH_NAME senzswitch
ENV SWITCH_HOST dev.localhost
ENV SWITCH_PORT 7070
ENV SENZIE_NAME sampath
ENV SENZIE_MODE DEV
ENV CASSANDRA_HOST dev.localhost
ENV CASSANDRA_PORT 9042

# fincale config
ENV FINACLE_API https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService
ENV LIEN_ADD_ACTION https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService
ENV LIEN_MOD_ACTION https://fin10env1.sampath.lk:15250/fiwebservice/FIWebService

# finacle integrator config
ENV TRANS_API http://192.168.125.93:7800/sd/iib/IIBFinacleIntegration
ENV TRANS_ACTION http://192.168.125.93:7800/sd/iib/iibfinacleintegration
ENV TRANS_ACCOUNT 900010000801
ENV TRANS_COMMISSION 20

# copy app
ADD . /app
WORKDIR /app

# build
RUN go build -o build/senz src/*.go

# server running port
EXPOSE 7070

# .keys volume
VOLUME ["/app/.keys"]

ENTRYPOINT ["/app/docker-entrypoint.sh"]
