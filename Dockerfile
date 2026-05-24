# STAICFILES BUILDER
FROM node:alpine as web-build

WORKDIR /app

COPY package.json .
COPY package-lock.json .
COPY assets/ ./assets

RUN npm i && npm run sass:build
RUN cp -r ./assets/*/**.svg ./static

# APPLICATION BUILDER
FROM golang:alpine as server-build

WORKDIR /app

COPY go.mod .
COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY --from=web-build /app/static/ ./cmd/server/static

RUN echo | ls -lar

RUN go build cmd/server/main.go

# APPLICTION
FROM alpine

ENV STUB_HOST=${STUB_HOST:-0.0.0.0}
ENV STUB_PORT=${STUB_PORT:-8080}
ENV STUB_TELEMT_HOST=${STUB_TELEMT_HOST:-localhost:9091}

WORKDIR /app

COPY --from=server-build /app/main ./
ENTRYPOINT $PWD/main