FROM golang:1.11-alpine3.8 as build
COPY . /go/src/app
RUN apk add --no-cache --update build-base git && \
    cd /go/src/app/ && \
    make

FROM alpine:3.8
RUN apk add --no-cache --update ca-certificates
COPY --from=build /go/src/app/bin/seattlewaste2mqtt /usr/local/bin/seattlewaste2mqtt
VOLUME /config
CMD seattlewaste2mqtt -c /config/config.yaml
