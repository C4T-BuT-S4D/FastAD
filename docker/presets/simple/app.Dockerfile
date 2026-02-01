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
            -o "/allinone" \
            "./cmd/allinone/main.go"

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
        go build \
            -trimpath \
            -ldflags="-w -s" \
            -o "/migrator" \
            "./cmd/migrator/main.go"

ARG COMPRESS_BINARIES="false"
ENV COMPRESS_BINARIES=$COMPRESS_BINARIES
RUN if [ "${COMPRESS_BINARIES}" = "true" ]; then upx --lzma -9 "/allinone" && upx --lzma -9 "/migrator"; fi

FROM python:3.12-bookworm

ENV PWNLIB_NOTERM=true
ENV UV_BREAK_SYSTEM_PACKAGES=true
ENV PIP_BREAK_SYSTEM_PACKAGES=true

COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/

ARG CHECKERS_DIR
COPY ${CHECKERS_DIR} /checkers
RUN --mount=type=cache,target=/root/.cache/uv \
    find /checkers -name "requirements.txt" | while read -r file; do \
      echo "Installing requirements from ${file}"; uv pip install --system -r "$file"; \
    done

COPY --from=builder /allinone /allinone
COPY --from=builder /migrator /migrator

CMD ["/bin/sh", "-c", "/migrator init && /migrator migrate && /allinone"]
