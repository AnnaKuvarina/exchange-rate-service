FROM golang:1.22.3-alpine3.19 as build-env

# Copy the source from the current directory to the Working Directory inside the container
WORKDIR /app

#Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Build the Go app (cache if not changes are made)
RUN --mount=type=cache,target=/root/.cache/go-build DOCKER_BUILDKIT=1 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix cgo -o bin/exchange-rate-service cmd/*


#runtime image
FROM alpine:3

USER root

COPY --from=build-env /app/bin/exchange-rate-service /app/bin/exchange-rate-service
COPY --from=build-env /app/db /app/db

# Run
ENTRYPOINT ["/app/bin/exchange-rate-service"]
