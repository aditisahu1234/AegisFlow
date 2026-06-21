# -----------------------------
# Stage 1 - Build the application
# -----------------------------

#never hardcore versions, CI/CD inject em automatically during builds
ARG VERSION=1.0.0
ARG BUILD_DATE
ARG VCS_REF

#starts from official Go image
#copy only executable from builder, not everything
FROM golang:1.26.4 AS builder 

#creates app and moves docker inside it, working directory
WORKDIR /app

#docker won't redownload go mods as they remian unchanged when
#source code changes, only copies dependency definition files
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
#build the go library, reproducible builds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o aegisflow .

# ==============================
# Stage 2 — Production Runtime
# ==============================

#production runtime image, contains only your binary, nothing else
#nothing an attacker can use
FROM gcr.io/distroless/static-debian12

ARG VERSION=1.0.0
ARG BUILD_DATE
ARG VCS_REF
#OCI (open Container Initiative) image specification
LABEL org.opencontainers.image.title="AegisFlow"

LABEL org.opencontainers.image.description="Adaptive Distributed API Protection Platform"

LABEL org.opencontainers.image.version=$VERSION
LABEL org.opencontainers.image.created=$BUILD_DATE
LABEL org.opencontainers.image.revision=$VCS_REF

LABEL org.opencontainers.image.authors="Aditi Sahu"

LABEL org.opencontainers.image.source="https://github.com/aditisahu1234/AegisFlow"

LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /

#only copy your binary, executable code for production
COPY --from=builder /app/aegisflow .
COPY --from=builder /app/redis/scripts /redis/scripts

EXPOSE 8080

#run as non root
USER nonroot:nonroot

ENTRYPOINT ["./aegisflow"]