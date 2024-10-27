# STEP 1: Build with golang image
FROM golang:1.23 AS builder
# Add Maintainer Info
LABEL maintainer="Bernardo Secades <bernardosecades@gmail.com>"

RUN useradd appuser

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/bernardosecades/quizapp
# Copy go mod and sum files
COPY ../../go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY ../.. .
# Build the Go app

# -ldflags="-w -s" reduce size of bnary
RUN cd cmd/cli && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/cli-quizapp .

# Command to run the executable
ENTRYPOINT ["/bin/cli-quizapp"]

# STEP 2: Build a small image
FROM alpine:3.13 AS production

# passwd file to appuser created in first stage
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/cli-quizapp /go/bin/cli-quizapp

RUN  ls -la /go/bin/

# Use an unprivileged user.
USER appuser

# Command to run the executable
ENTRYPOINT ["/go/bin/cli-quizapp"]