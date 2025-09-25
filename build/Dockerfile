FROM golang:1.24 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o docs_storage_service ./cmd/app/

FROM alpine:latest AS build-release-stage

WORKDIR /docs_storage_service

COPY --from=build-stage /app/docs_storage_service .

CMD ["./docs_storage_service"]
