# -----------------------------
# Stage 1 - Build the application
# -----------------------------

#starts from official Go image
#copy only executable from builder, not everything
FROM golang:1.26.4 AS builder 

#creates app and moves docker inside it, working directory
WORKDIR /app

#docker won't redownload go mods as they remian unchanged when
#source code changes, only copies dependency definition files
COPY go.mod go.sum ./

RUN go mod download

COPY . .
#build the go library
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o aegisflow .

# ==============================
# Stage 2 — Production Runtime
# ==============================

#production runtime image, contains only your binary, nothing else
#nothing an attacker can use
FROM gcr.io/distroless/static-debian12

WORKDIR /

#only copy your binary, executable code for production
COPY --from=builder /app/aegisflow .

EXPOSE 8080

#run as non root
USER nonroot:nonroot

ENTRYPOINT ["./aegisflow"]