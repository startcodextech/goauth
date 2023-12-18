FROM golang:latest as dev

WORKDIR /goauth
COPY go.mod ./
COPY . .

RUN go mod download

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./server ./cmd/server

FROM alpine:latest as system

RUN apk --no-cache add ca-certificates wget tzdata && update-ca-certificates

ENV TZ=UTC

COPY --from=dev /goauth/server /goauth/server

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
  CMD [ "wget", "localhost:8000/health", "-q", "-O", "-" ]

EXPOSE 8000
EXPOSE 8080

CMD ["/goauth/server"]