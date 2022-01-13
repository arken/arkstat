# Start from the latest golang base image
FROM ghcr.io/autamus/go:latest as builder

ARG version

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -ldflags "-s -w -X github.com/arken/arkstat/config.Version=$version" -o arkstat .

# Start again with minimal envoirnment.
FROM ubuntu:latest

RUN apt update && \
    apt install -y ca-certificates wget

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/arkstat /app/arkstat
COPY mail/templates /app/mail/templates
COPY entrypoint.sh  /app/entrypoint.sh

RUN wget https://github.com/benbjohnson/litestream/releases/download/v0.3.7/litestream-v0.3.7-linux-amd64.deb && \
    dpkg -i litestream-v0.3.7-linux-amd64.deb

RUN chmod a+x /app/entrypoint.sh

# Command to run the executable
ENTRYPOINT ["/app/entrypoint.sh"]
