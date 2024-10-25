FROM golang:1.22.3-alpine as build-price-api
LABEL version="0.1.0"

ARG GITLAB_USERNAME
ARG GITLAB_PASSWORD
ARG GO_MOD_TAG

RUN apk add jq make git

ENV GOPRIVATE=github.com/ipanardian/price-api

RUN echo "machine gitlab.com login $GITLAB_USERNAME password $GITLAB_PASSWORD" > ~/.netrc

WORKDIR /etc/price-api/

COPY . .

RUN make get-module GO_MOD_TAG=$GO_MOD_TAG

RUN make build-linux

FROM bash AS price-api

RUN apk add supervisor htop busybox-extras vim

WORKDIR /etc/price-api/

COPY --from=build-price-api /etc/price-api/build .
COPY --from=build-price-api /etc/price-api/build /usr/bin

COPY ./supervisord/api.conf /etc/supervisor/conf.d/
COPY ./supervisord/consumer.conf /etc/supervisor/conf.d/
COPY ./supervisord/supervisord.conf /etc/supervisor/

COPY .env .env

ENTRYPOINT ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]