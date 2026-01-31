FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache upx

WORKDIR /app

COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY go.* ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
        go build \
            -trimpath \
            -ldflags="-w -s" \
            -o "/checkers" \
            "./cmd/checkers/main.go"

ARG COMPRESS_BINARIES="false"
ENV COMPRESS_BINARIES=$COMPRESS_BINARIES
RUN if [ "${COMPRESS_BINARIES}" = "true" ]; then upx --lzma -9 "/checkers"; fi

FROM python:3.12-bookworm

ENV PWNLIB_NOTERM=true
ENV UV_BREAK_SYSTEM_PACKAGES=true
ENV PIP_BREAK_SYSTEM_PACKAGES=true

COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/

ARG CHECKERS_DIR
COPY ${CHECKERS_DIR} /checkers-scripts
RUN --mount=type=cache,target=/root/.cache/uv \
    find /checkers-scripts -name "requirements.txt" | while read -r file; do \
      echo "Installing requirements from ${file}"; uv pip install --system -r "$file"; \
    done

COPY --from=builder /checkers /checkers

ENTRYPOINT ["/checkers"]
