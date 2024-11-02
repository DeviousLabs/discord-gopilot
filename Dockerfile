# Stage 1: Build Stage
FROM alpine:edge AS builder

# Install necessary packages
RUN apk update && apk add --no-cache \
    bash \
    git \
    go \
    && rm -rf /var/cache/apk/*

# Set environment variables for Go
ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Prepare the working directory for the Go app
RUN mkdir -p "$GOPATH/src/gopilot" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Copy the Go application source into the image
COPY . $GOPATH/src/gopilot/

# Install Garble
RUN go install mvdan.cc/garble@latest

# Set the working directory to your app's location
WORKDIR $GOPATH/src/gopilot

# Use garble to build the Go application
RUN PATH=$(go env GOROOT)/bin:${PATH} garble -literals -tiny build -o /usr/local/bin/gopilot ./cmd

# Stage 2: Final Image
FROM alpine:latest

# Install runtime dependencies only
RUN apk add --no-cache \
    bash \
    && rm -rf /var/cache/apk/*

# Copy the built executable from the builder stage
COPY --from=builder /usr/local/bin/gopilot /usr/local/bin/gopilot

# Set permissions and create user
RUN chmod +x /usr/local/bin/gopilot \
    && addgroup -S gopilot && adduser -S gopilot -G gopilot
USER gopilot

# Set up environment for runtime
WORKDIR /home/gopilot
VOLUME /home/gopilot/data

CMD ["/usr/local/bin/gopilot"]