FROM golang:1.23-alpine3.20 as builder

RUN apk update && apk add --no-cache ca-certificates

# Install libwebp-tools for cwebp and dwebp binaries and ffmpeg for video processing
RUN apk add --no-cache libwebp-tools libpng libjpeg-turbo giflib tiff ffmpeg

# Create the .bin/webp directory and link cwebp binary for go-webpbin
RUN mkdir -p /root/.bin/webp && \
    ln -s /usr/bin/cwebp /root/.bin/webp/cwebp

WORKDIR /app
COPY . .

COPY ./src/docs/swagger.yaml /app/src/docs/swagger.yaml

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /main .

FROM scratch
WORKDIR /
COPY --from=builder /main /main
COPY --from=builder /app/src/docs /src/docs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/main"]

