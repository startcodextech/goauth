FROM golang:latest as dev

WORKDIR /goauth
COPY go.mod ./
COPY . .

RUN go mod download

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./service ./cmd/service

FROM alpine:latest as system

COPY --from=dev /goauth/service /goauth/server

RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN apk --no-cache add wget

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
  CMD [ "wget", "localhost:8080/health", "-q", "-O", "-" ]

EXPOSE 8080

CMD ["/goauth/server"]