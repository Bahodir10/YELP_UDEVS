FROM golang:1.23.3-alpine

RUN apk update && apk add bash postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY internal /app/internal
COPY . .

RUN go build -o main ./cmd/server/main.go

EXPOSE 9090

ENTRYPOINT ["sh", "-c", "until pg_isready -h db -p 5432; do echo waiting for db; sleep 2; done; ./main"]
