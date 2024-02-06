 

FROM node:alpine as react-builder

WORKDIR /app

COPY ui/package.json ui/yarn.lock ./

RUN yarn install --frozen-lockfile

COPY ui/ ./

RUN yarn build


FROM golang:alpine as go-builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .git ./.git
COPY cmd ./cmd
COPY pkg ./pkg
COPY version.go ./version.go
COPY ui/ui.go ./ui/ui.go
COPY Makefile ./Makefile

COPY --from=react-builder /app/build ./ui/build

ARG VERSION=dev

RUN apk add --no-cache git make

RUN make build-go

FROM alpine:latest

RUN apk upgrade --no-cache && apk --no-cache add ca-certificates

COPY --from=go-builder /app/bin/aartiserver /usr/local/bin/aartiserver
COPY --from=go-builder /app/bin/aarticlient /usr/local/bin/aarticlient

ENTRYPOINT ["/usr/local/bin/aartiserver"]
