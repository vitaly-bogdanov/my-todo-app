FROM node:22-bookworm AS web
WORKDIR /src
COPY web/package.json web/package-lock.json web/
RUN npm --prefix web ci
COPY web web
RUN npm --prefix web run build

FROM golang:1.26.4-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/todoapp/exe /app/cmd/todoapp

FROM alpine:3.23
# tzdata нужен для TIME_ZONE, отличных от UTC; ca-certificates — для sslmode=require/verify-*.
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
RUN mkdir -p /app/out/logs
COPY --from=builder /app/cmd/todoapp/exe /app/cmd/todoapp/exe
COPY --from=web /src/web/dist /app/web/dist
EXPOSE 8080 10000
CMD ["/app/cmd/todoapp/exe"]
