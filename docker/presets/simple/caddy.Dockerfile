FROM node:22.14.0-alpine AS front-base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

COPY front /app
WORKDIR /app
RUN corepack enable

FROM front-base AS front-build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM caddy:2.9.1-alpine

COPY --from=front-build /app/dist /front
