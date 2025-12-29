FROM golang:1.25-bookworm AS builder

ARG versionflags

WORKDIR /app

RUN apt-get update && apt-get upgrade -y && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -a -tags netgo -ldflags="-extldflags '-static' -s -w $versionflags" -o wg cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get upgrade -y && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        wget \
        dnsutils \
        iputils-ping \
        vim \
        nano \
        jq \
        lsof \
        net-tools \
        procps \
        tzdata \
        traceroute \
        mtr-tiny \
        && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/wg /app/wg

RUN chmod +x /app/wg

ENV PATH="/app:${PATH}"
