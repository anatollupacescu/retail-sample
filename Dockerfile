# linter
FROM golang:1.13

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

ADD go.mod go.sum /m/
RUN cd /m && go mod download

RUN mkdir -p /retail
ADD . /retail
WORKDIR /retail

# Run linters
RUN golangci-lint run --issues-exit-code=1 --deadline=600s ./...

# Run tests
RUN go test -timeout=600s -v --race ./...

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 make build
