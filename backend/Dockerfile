FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/bin/api /app/bin/api
EXPOSE 8443
CMD ["/app/bin/api"]
