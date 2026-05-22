ARG APP_NAME=app

FROM node:alpine as web-build

WORKDIR /app

COPY package.json .
COPY package-lock.json .
COPY assets/ ./assets

RUN npm i && npm run sass:build
RUN cp -r ./assets/*/**.svg ./static

FROM golang:alpine as server-build
ARG APP_NAME


WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY templates ./templates
COPY --from=web-build /app/static/ ./static

RUN go build -o ${APP_NAME}

# APPLICTION
FROM alpine

ARG APP_NAME
ENV APP_NAME=${APP_NAME}
ENV STUB_PORT=${STUB_PORT:-8080}
ENV STUB_TELEMT_LISTEN=${STUB_TELEMT_LISTEN:-localhost:9091}

WORKDIR /app

COPY --from=server-build /app/app ./
ENTRYPOINT $PWD/$APP_NAME