FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN apk add --update make gcc musl-dev
RUN make build

FROM alpine:latest
WORKDIR /app
RUN apk add --update ffmpeg
COPY --from=builder /build/bin/hola /app

CMD ["./hola"]