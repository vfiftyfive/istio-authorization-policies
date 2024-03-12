# Start from the official Go image to create a build artifact.
FROM golang:alpine AS builder

# Install git, required for fetching Go dependencies.
RUN apk update && apk add --no-cache git

# Set the working directory inside the container.
WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . .

# Build the command inside the container.
# (CGO_ENABLED=0 for a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest

# Copy the binary from the builder stage.
COPY --from=builder /app/main .

# Make port 8081 available to the outside world.
EXPOSE 8081

# Run the binary.
CMD ["./main"]