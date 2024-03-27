# Official Go base image
FROM golang:latest

# Working dir inside container
WORKDIR /go/src/app

# Copy local package files into container worskpace
COPY . .

# Download dependencies
RUN go mod download

# Build Go app
RUN go build -o dicom-service cmd/api-server/main.go

# Expose port for Go app
EXPOSE 8080

# Command to run the exec
CMD ["./dicom-service"]
