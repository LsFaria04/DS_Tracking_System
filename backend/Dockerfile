# First select the go version that we are going to work with
FROM golang:1.25.3 As builder

# Select the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY src/go.* ./

# Download all dependencies
RUN go mod download

# Copy the source files to the container
COPY src/ .

# Build the app
RUN go build -o app .

# Stage 2: minimal runtime image
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
# Run the binary
CMD ["./app"]