FROM --platform=$BUILDPLATFORM golang:1.25.3 AS build

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod ./
COPY go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd ./cmd
COPY pkg ./pkg

# Automatically provided by the buildkit
ARG TARGETOS
ARG TARGETARCH

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -a -o lazygate cmd/lazygate/main.go

# Move binary into final image
FROM --platform=$BUILDPLATFORM gcr.io/distroless/static-debian11 AS app
COPY --from=build /workspace/lazygate /gate
CMD ["/gate"]
