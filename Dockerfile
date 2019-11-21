ARG GO_VERSION=1.13

FROM golang:${GO_VERSION}-alpine AS builder

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates git
RUN apk update && apk upgrade && apk add bash
# Set the environment variables for the go command:
# * CGO_ENABLED=0 to build a statically-linked executable
# * GOFLAGS=-mod=vendor to force `go build` to look into the `/vendor` folder.
ENV CGO_ENABLED=0

## Prebuild requirements
RUN go get golang.org/x/tools/cmd/goimports
RUN wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# test before build
#RUN test -z "$(goimports -d $(find . -type f -name '*.go' -not -path "./vendor/*"|grep -v -f .goimportsignore) 2>&1 | tee /dev/fd/2)"
#RUN test -z "$(find . -type f -name '*.go'|grep -v -f .goimportsignore|xargs goimports -d 2>&1 |xargs grep -nE '\serrors.New\('| grep -vE '^[^[:space:]]+:\s+(var\s+)?Err[^[:space:]]+\s*=\s*errors.New\(' 2>&1 | tee /dev/fd/2)"
#RUN golangci-lint run
#RUN go test -v ./...

# Build the executable to `/app`. Mark the build as statically linked.
RUN go build ./src/cmd/queue

# Final stage: the running container.
FROM scratch AS final

WORKDIR /app

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the second stage.
COPY --from=builder /src/queue /app/queue

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080

# Run the compiled binary.
CMD ["/app/queue"]
