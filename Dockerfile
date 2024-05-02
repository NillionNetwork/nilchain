FROM golang:1.22-alpine3.18 as builder

RUN set -eux; apk add --no-cache git libusb-dev linux-headers gcc musl-dev make;

ENV GOPATH=""

# Copy relevant files before go mod download. Replace directives to local paths break if local
# files are not copied before go mod download.
ADD app app
ADD nilliond nilliond
ADD params params

COPY Makefile .
COPY go.mod .
COPY go.sum .

RUN go mod download

RUN make build

FROM alpine:3.18

COPY --from=builder /go/build/nilliond /bin/nilliond

ENTRYPOINT ["nilliond"]
