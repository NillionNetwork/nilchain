FROM golang:1.22-alpine3.18 as builder
ARG BINARY_VERSION

RUN set -eux; apk add --no-cache git libusb-dev linux-headers gcc musl-dev make;

ENV GOPATH=""

# Copy relevant files before go mod download. Replace directives to local paths break if local
# files are not copied before go mod download.
ADD app app
ADD nilchaind nilchaind
ADD params params
ADD x x 

ADD .git .git

COPY Makefile .
COPY go.mod .
COPY go.sum .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false go build -mod=readonly -tags "netgo ledger" -ldflags '-X github.com/cosmos/cosmos-sdk/version.Name=sim -X github.com/cosmos/cosmos-sdk/version.AppName=simd -X github.com/cosmos/cosmos-sdk/version.Version= -X github.com/cosmos/cosmos-sdk/version.Commit=$BINARY_VERSION -X "github.com/cosmos/cosmos-sdk/version.BuildTags=netgo ledger," -w -s' -trimpath -o /go/build/ ./...

FROM alpine:3.18

COPY --from=builder /go/build/nilchaind /bin/nilchaind

ENTRYPOINT ["nilchaind"]
