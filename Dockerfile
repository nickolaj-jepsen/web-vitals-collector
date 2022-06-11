FROM golang:1.18-alpine as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server cmd/main/main.go
RUN go build -v -o migrate cmd/migrate/main.go

FROM alpine

WORKDIR /

COPY migrations /migrations
COPY --from=builder /app/server /server
COPY --from=builder /app/migrate /migrate

CMD /migrate ; /server
