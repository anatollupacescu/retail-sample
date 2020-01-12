# linter
FROM golang:1.13 as tester

# golangci-lint version
ENV VERSION 1.18.0
ENV CHECKSUM 0ef2c502035d5f12d6d3a30a7c4469cfcae4dd3828d15fbbfb799c8331cd51c4

# store checksum in a file to be able to compare against it
RUN echo "${CHECKSUM}  golangci-lint-${VERSION}-linux-amd64.tar.gz" > CHECKSUM

# Download from Github the specified release and extract into the go/bin folder
RUN curl -L "https://github.com/golangci/golangci-lint/releases/download/v${VERSION}/golangci-lint-${VERSION}-linux-amd64.tar.gz" \
    -o golangci-lint-${VERSION}-linux-amd64.tar.gz \
    && shasum -a 256 -c CHECKSUM \
    && tar xvzf golangci-lint-${VERSION}-linux-amd64.tar.gz \
    --strip-components=1 \
    -C ./bin \
    golangci-lint-${VERSION}-linux-amd64/golangci-lint

# clean up
RUN rm -rf CHECKSUM "golangci-lint-${VERSION}-linux-amd64.tar.gz"

RUN mkdir -p /retail
ADD . /retail
WORKDIR /retail

# Run linters
RUN golangci-lint run --issues-exit-code=1 --deadline=600s ./...

# Run tests
RUN go test -timeout=600s -v --race ./...

FROM golang:1.13 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.13 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /retail
ADD . /retail
WORKDIR /retail

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 make build

# Final stage: Run the binary
FROM scratch

ENV PORT 8080
ENV DIAG_PORT 8181

# certificates to interact with other services
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER myapp

# and finally the binary
COPY --from=builder /retail/bin/retail /retail

EXPOSE $PORT
EXPOSE $DIAG_PORT

ADD web /web

CMD ["/retail"]
