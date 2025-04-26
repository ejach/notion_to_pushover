FROM golang:1.24.2-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o app

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/app .

EXPOSE 8069

ENTRYPOINT ["./app"]