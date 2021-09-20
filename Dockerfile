FROM golang:1.17-alpine

WORKDIR /accountclient
ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENTRYPOINT ["/accountclient/entrypoint.sh"]
