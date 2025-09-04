FROM golang:latest as builder

RUN apt-get update && apt-get install -y \
  make \
  build-essential \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /anvil-gateway

COPY bin/gateway .
COPY Makefile .

RUN adduser -D -g "" noo8xl
USER noo8xl

ENV GO_ENV=production

EXPOSE 30201

FROM golang:latest as runner



CMD ["./bin/gateway"]
