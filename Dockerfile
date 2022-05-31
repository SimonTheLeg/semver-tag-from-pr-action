# syntax=docker/dockerfile:1.4
FROM golang:1.18.2-alpine3.16 AS alpine-upx

LABEL org.opencontainers.image.source https://github.com/simontheleg/semver-tag-from-pr-action

RUN apk update && apk add upx binutils

FROM alpine-upx AS builder

ARG outpath="/bin/action"
WORKDIR /build

COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go


# TODO cacheing does not seem to work properly. Needs to be inspected
RUN --mount=type=cache,target=root/.cache/go-build \
  CGO_ENABLED=0 \
  go build \
  -trimpath \
  -ldflags '-w -s' \
  -o ${outpath}


# Strip any symbols
RUN strip ${outpath}
# Compress the compiled action
RUN upx -q -9 ${outpath}

FROM scratch
ARG outpath="/bin/action"
COPY --from=builder ${outpath} ${outpath}

ENTRYPOINT [ "/bin/action" ]