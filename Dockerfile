FROM golang:1.23.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -o /trackingapi cmd/main.go


FROM scratch AS build-release-stage

WORKDIR /

COPY --from=builder /trackingapi /trackingapi

EXPOSE 8080

ENTRYPOINT ["/trackingapi"]