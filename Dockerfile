FROM golang:latest as base

WORKDIR /goauth

COPY go.mod ./
COPY . .

RUN go mod download

RUN go mod tidy

EXPOSE 80
EXPOSE 9090

FROM base as dev

WORKDIR /goauth

RUN go install github.com/cosmtrek/air@latest

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./tmp/main ./cmd/server/main.go

CMD ["air", "-c", ".air.toml"]

FROM base as build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./server ./cmd/server

FROM alpine:latest as system

RUN apk --no-cache add ca-certificates wget tzdata && update-ca-certificates

ENV TZ=UTC

COPY --from=build /goauth/server /goauth/server

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
  CMD [ "wget", "localhost:8000/health", "-q", "-O", "-" ]

EXPOSE 80
EXPOSE 9090

CMD ["/goauth/server"]