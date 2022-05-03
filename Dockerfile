FROM golang:alpine AS build
RUN apk add git gcc musl-dev
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 1
ADD . /go/src/Vanilla/
WORKDIR /go/src/Vanilla
RUN go build .

FROM alpine
COPY --from=build /go/src/Vanilla/Vanilla /bin/Vanilla
WORKDIR /data
CMD Vanilla