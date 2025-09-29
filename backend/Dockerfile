# ---------- Builder ----------
FROM golang:1.23-alpine3.20 AS builder

RUN apk update && apk add --no-cache ca-certificates \
  libwebp-tools libpng libjpeg-turbo giflib tiff ffmpeg

# เตรียม path สำหรับ go-webpbin (ตอน build เฉยๆ ยังไม่จำเป็น แต่ไม่เป็นไร)
RUN mkdir -p /root/.bin/webp && ln -s /usr/bin/cwebp /root/.bin/webp/cwebp

# Install libwebp-tools for cwebp and dwebp binaries and ffmpeg for video processing
RUN apk add --no-cache libwebp-tools libpng libjpeg-turbo giflib tiff ffmpeg

# Create the .bin/webp directory and link cwebp binary for go-webpbin
RUN mkdir -p /root/.bin/webp && \
    ln -s /usr/bin/cwebp /root/.bin/webp/cwebp

WORKDIR /app
COPY . .
COPY ./src/docs/swagger.yaml /app/src/docs/swagger.yaml

# สร้างไบนารีแบบ static
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main .

# ---------- Final (Alpine) ----------
FROM alpine:3.20

# ติดตั้ง runtime deps สำหรับ cwebp/ffmpeg และ certs
RUN apk add --no-cache ca-certificates \
  libwebp-tools libpng libjpeg-turbo giflib tiff ffmpeg

WORKDIR /app

# คัดลอกไบนารีและ docs
COPY --from=builder /app/main /app/main
COPY --from=builder /app/src/docs/ /app/src/docs/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# สร้าง path ที่ go-webpbin เรียกหา: .bin/webp/cwebp (relative จาก WORKDIR)
RUN mkdir -p /app/.bin/webp && ln -s /usr/bin/cwebp /app/.bin/webp/cwebp

# รัน
ENTRYPOINT ["/app/main"]
