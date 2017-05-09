FROM golang:1.8-alpine

# Add files
ADD . /go/src/github.com/wattx/gke-snapshooter
WORKDIR /go/src/github.com/wattx/gke-snapshooter

# Install the project
RUN go install .

# Entrypoint
ENTRYPOINT ["/go/bin/gke-snapshooter"]
CMD ["--help"]
