# syntax=docker/dockerfile:1.4

FROM golang:1.18.2-alpine3.16 AS alpine-upx

RUN apk update && apk add upx binutils

FROM alpine-upx AS builder

ARG outpath="/bin/action"
WORKDIR /build

COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go


RUN --mount=type=cache,target=root/.cache/go-build \
  go build -o ${outpath}


# Strip any symbols
RUN strip ${outpath}
# Compress the compiled action
RUN upx -q -9 ${outpath}

FROM scratch
ARG outpath="/bin/action"
COPY --from=builder ${outpath} ${outpath}

ENTRYPOINT [ "/bin/action" ]