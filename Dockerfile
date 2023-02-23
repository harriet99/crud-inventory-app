FROM golang:1.17

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download all dependencies
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the application inside the container
RUN go build -o main .

# Expose port 8080 for the application to listen on
EXPOSE 8080

# Run the application when the container starts
CMD ["./main"]






