FROM golang:alpine AS build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install packages required by the image
RUN apk update && \
    apk add curl \
    git \
    bash \
    make \
    ca-certificates && \
    rm -rf /var/cache/apk/*


WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . .
RUN make build

COPY /scripts/entrypoint.sh .
RUN ls -la
ENTRYPOINT ["bash", "./entrypoint.sh"]
