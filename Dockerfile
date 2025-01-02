FROM golang:1.23.4-alpine AS builder
WORKDIR /app

RUN apk add --no-cache openssl

RUN openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -o /trackingapi cmd/api/main.go


FROM scratch AS build-release-stage

WORKDIR /

COPY --from=builder /trackingapi /trackingapi

COPY --from=builder /app/server.crt /server.crt
COPY --from=builder /app/server.key /server.key

EXPOSE 8080

ENTRYPOINT ["/trackingapi"]