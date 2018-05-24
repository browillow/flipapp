# 1) Build the go web server

FROM golang:1.9-alpine as build-go
RUN apk --no-cache add git
WORKDIR /go/src/browillow/flipper/flipperapp
COPY main.go /go/src/browillow/flipper/flipperapp
RUN go install

# 2) Build the Angular web app

FROM node:8.9.1-alpine AS build-Angular
RUN apk --no-cache add python2
# This is required due to this issue: https://github.com/nodejs/node-gyp/issues/1236#issuecomment-309401410
RUN mkdir /root/.npm-global && npm config set prefix '/root/.npm-global'
ENV PATH="/root/.npm-global/bin:${PATH}"
ENV NPM_CONFIG_LOGLEVEL warn
ENV NPM_CONFIG_PREFIX=/root/.npm-global
RUN npm install -g npm@latest
RUN mkdir -p /build
COPY angular.json /build/
COPY package-lock.json /build/
COPY package.json /build/
COPY tsconfig.json /build/
COPY tslint.json /build/
COPY /src /build/src
COPY /e2e /build/e2e
RUN cd /build && npm install
RUN cd /build && npm run build

# 3) Build the final image

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app/server/
COPY --from=build-go /go/bin /app/server
COPY --from=build-angular /build/dist /app/server/dist
ENTRYPOINT /app/server/flipperapp
EXPOSE 80
