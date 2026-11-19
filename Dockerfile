# syntax=docker/dockerfile:1
# rancher-projects image: harbor.support.tools/rancher-projects/rancher-projects (cluster-services-1489).
#
# Fetch tiers (no package-manager fetch anywhere in this build):
#   base image  golang:1.27.4-trixie through the Harbor Docker Hub proxy-cache, pinned by index digest.
#               It already ships openssl, awk, find and ca-certificates (baked layers), which is all
#               noemsign-ca.sh needs, so the build runs no apk/apt.
#   Go modules  Athens (GOPROXY=https://goproxy.support.tools, a cache we run); checksum DB stays on.
#
# CA standard (cluster-services-1375): noemsign-ca.sh (vendored, identical sha256 to
# SupportTools/cluster-services noemsign/noemsign-ca.sh) strips every emSign root and its guard fails
# the build if any is left. The cleaned bundle is the only trust store in the scratch runtime image.
#
# The builder runs on the build host's platform and cross-compiles for TARGETARCH, so the multi-arch
# build needs no QEMU and the runtime stage has no RUN step.
FROM --platform=$BUILDPLATFORM harbor.support.tools/dockerhub/library/golang:1.27.4-trixie@sha256:f18c3d73a184fb0582828510722e723b63a93f7188d11dd223d9be6f6880e65d AS builder

RUN --mount=type=bind,source=noemsign-ca.sh,target=/tmp/noemsign-ca.sh \
    sh /tmp/noemsign-ca.sh strip / \
 && sh /tmp/noemsign-ca.sh guard /

WORKDIR /src
ARG GOPROXY=https://goproxy.support.tools
ENV GOPROXY=${GOPROXY} GOTOOLCHAIN=local
COPY go.mod go.sum ./
# -x prints every "# get <url>" line, so the build log shows which proxy served each module.
RUN echo "go module source: GOPROXY=${GOPROXY}" && go mod download -x

COPY . .

ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w \
    -X github.com/supporttools/rancher-projects/pkg/version.Version=${VERSION} \
    -X github.com/supporttools/rancher-projects/pkg/version.GitCommit=${GIT_COMMIT} \
    -X github.com/supporttools/rancher-projects/pkg/version.BuildTime=${BUILD_DATE}" \
    -o /out/rancher-projects

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/rancher-projects /rancher-projects
USER 65534:65534
ENTRYPOINT ["/rancher-projects"]
CMD ["--help"]
