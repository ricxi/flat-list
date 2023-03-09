FROM golang:1.19-alpine3.17 AS base
WORKDIR /app
ENV GOPRIVATE=github.com/ricxi/flat-list
COPY . .
RUN go mod download
RUN go build -o bin/userService ./cmd/http

FROM alpine:3.17 AS dev-build
WORKDIR /
COPY --from=base /app/bin .
# ENV PORT=80
# EXPOSE 80
CMD ./userService