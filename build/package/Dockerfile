FROM golang:1.14 as build-env

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/hello

FROM scratch

RUN apk add --update --no-cache ca-certificates git

COPY --from=build-env /go/bin/hello /go/bin/hello

ENTRYPOINT ["/go/bin/hello"]