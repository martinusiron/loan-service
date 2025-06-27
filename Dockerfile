FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init -g cmd/main.go --parseDependency --parseInternal

RUN go build -v -o app ./cmd/main.go && ls -l ./app

EXPOSE 8080

CMD ["./app"]
