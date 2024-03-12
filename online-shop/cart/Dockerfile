# Start with the Go image to create a build artifact.
FROM golang:1.22 AS builder

# Set the work directory inside the container.
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the local package files to the container's workspace.
COPY *.go ./

# Build the application inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a multi-stage build to create a lean image.
FROM alpine:latest  

# Copy the binary from the builder stage.
COPY --from=builder /app/main .

# Expose port 8080.
EXPOSE 8083

# Run the binary.
CMD ["./main"]
