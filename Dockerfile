FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
RUN echo "running on $BUILDPLATFORM, building for $TARGETPLATFORM - $TARGETOS $TARGETARCH"
RUN apk update && apk add --no-cache make git openssh ca-certificates
RUN mkdir /build
ADD . /build/
WORKDIR /build
ENV CGO_ENABLED=0
ENV GOARCH=$TARGETARCH
ENV GOOS=$TARGETOS
RUN go build -mod=vendor -ldflags="-s -w" -o devin ./cmd/server

# copy to scratch
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/devin /app/
WORKDIR /app
CMD ["./devin"]
