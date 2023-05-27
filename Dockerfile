FROM golang:1.20-alpine

WORKDIR /opt

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

COPY cmd/ ./cmd/
COPY req/ ./req/
COPY database/ ./database/

RUN go build -o ./pricebot

EXPOSE 8080

CMD [ "/opt/pricebot" ]