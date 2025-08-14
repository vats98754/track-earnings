# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/server /server
USER nonroot:nonroot
ENV ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["/server"]
