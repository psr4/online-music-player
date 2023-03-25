FROM golang:1.20

WORKDIR /service
#COPY go.mod go.sum ./
COPY . .

RUN  apt-get update \
    && apt-get install -y libtagc0-dev

ENTRYPOINT ["/bin/bash", "-c"]