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
WORKDIR /app
COPY --from=builder /app/cmd/todoapp/exe /app/cmd/todoapp/exe
COPY --from=web /src/web/dist /app/web/dist
CMD ["/app/cmd/todoapp/exe"]
