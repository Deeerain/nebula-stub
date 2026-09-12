# STAICFILES BUILDER
FROM node:alpine as web-build

WORKDIR /app

COPY frontend/package.json /app/
COPY frontend/package-lock.json /app/
COPY frontend/vite.config.js /app/
COPY frontend/public/ /app/public/
COPY frontend/src/ /app/src
COPY frontend/index.html /app/

RUN npm i && npm run build

# APPLICATION BUILDER
FROM golang:alpine as server-build

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY internal/ ./internal/
COPY main.go .
COPY --from=web-build /app/dist /app/frontend/dist

RUN echo | ls -lar
RUN echo | ls -lar ./internal

RUN go build main.go

# APPLICTION
FROM alpine

ENV STUB_HOST=${STUB_HOST:-0.0.0.0}
ENV STUB_PORT=${STUB_PORT:-8080}
ENV STUB_TELEMT_HOST=${STUB_TELEMT_HOST:-localhost:9091}

WORKDIR /app

COPY --from=server-build /app/main ./
ENTRYPOINT $PWD/main