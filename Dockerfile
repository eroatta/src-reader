# Set golang base image
FROM golang:alpine as builder

LABEL maintainer="Emiliano Roatta <emilianoroatta@gmail.com>"

# Install git to fetch the dependencies
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy dependency management files
COPY go.mod go.sum ./

# Fetch dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Start new stage: deployer
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage. Copy also the .env file.
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the main executable
CMD ["./main"]
