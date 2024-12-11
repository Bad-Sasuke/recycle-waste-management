FROM oven/bun:alpine as builder

WORKDIR /app

ARG VITE_WEB_API
ARG VITE_GITHUB_CLIENT_ID

COPY . /app

RUN bun install

RUN bun run build

FROM nginx:1.27.3-alpine3.20

## Copy our default nginx config
COPY nginx/default.conf /etc/nginx/conf.d/

## Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

## From ‘builder’ stage copy over the artifacts in dist folder to default nginx public folder
COPY --from=builder /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]