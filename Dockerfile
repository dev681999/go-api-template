FROM golang:1.15-alpine as backend_builder

RUN apk add upx

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o server -ldflags="-s -w" ./cmd/server/main.go
RUN upx -q -5 server

FROM alpine as prod

WORKDIR /app
COPY --from=backend_builder /app/server /

EXPOSE 8080
CMD [ "/server" ]
