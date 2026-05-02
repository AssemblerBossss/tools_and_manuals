FROM golang:1.19-alpine AS build

RUN apk add git

ENV GOBIN=/go/bin
WORKDIR /go/src/app
COPY *.mod ./
RUN go get main
COPY *.go ./

RUN go get -d -v ./...
RUN go build -v -o app ./

FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/app .
ENTRYPOINT ["./app"]
