# Build Stage
FROM golang:1.21-alpine3.18 AS builder

# Install necessary packages for building (including CGO support)
RUN apk add --no-cache bash build-base git

# Set the working directory
WORKDIR /app

# Copy the entire application directory into the container
COPY . .

# Build the Go application with CGO enabled
RUN CGO_ENABLED=1 go build -o main ./main.go

# Run Stage
FROM alpine:3.18

# Set the working directory in the second stage
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy other required assets like templates and static files
COPY templates templates/
COPY static static/

# Expose the port your application will listen on
EXPOSE 8080

# Define the command to run your application
CMD ["./main"]
