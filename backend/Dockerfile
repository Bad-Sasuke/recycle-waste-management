# ---------- Builder ----------
FROM golang:1.23-alpine3.20 AS builder

RUN apk update && apk add --no-cache ca-certificates \
  libwebp-tools libpng libjpeg-turbo giflib tiff ffmpeg

# ทำครั้งเดียวพอ ไม่ต้องมีซ้ำภายหลัง
RUN mkdir -p /root/.bin/webp && ln -s /usr/bin/cwebp /root/.bin/webp/cwebp

WORKDIR /app
COPY . .
COPY ./src/docs/swagger.yaml /app/src/docs/swagger.yaml
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main .

# ---------- Final ----------
FROM alpine:3.20
RUN apk add --no-cache ca-certificates libwebp-tools libpng libjpeg-turbo giflib tiff ffmpeg
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/src/docs/ /app/src/docs/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN mkdir -p /app/.bin/webp && ln -s /usr/bin/cwebp /app/.bin/webp/cwebp
ENTRYPOINT ["/app/main"]
