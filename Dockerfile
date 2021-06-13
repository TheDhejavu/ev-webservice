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


FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
RUN mkdir -p /var/log/app
WORKDIR /app/
COPY --from=build /app/server .
COPY --from=build /app/scripts/entrypoint.sh .
COPY --from=build /app/config/*.yml ./config/
RUN ls -la
ENTRYPOINT ["bash", "./entrypoint.sh"]
