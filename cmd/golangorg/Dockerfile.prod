# Builder
#########

FROM golang:1.12 AS build

RUN apt-get update && apt-get install -y \
      zip # required for generate-index.bash

# Check out the desired version of Go, both to build the golangorg binary and serve
# as the goroot for content serving.
ARG GO_REF
RUN test -n "$GO_REF" # GO_REF is required.
RUN git clone --single-branch --depth=1 -b $GO_REF https://go.googlesource.com/go /goroot
RUN cd /goroot/src && ./make.bash

ENV GOROOT /goroot
ENV PATH=/goroot/bin:$PATH
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

RUN go version

COPY . /website

WORKDIR /website/cmd/golangorg
RUN GOLANGORG_DOCSET=/goroot ./generate-index.bash

RUN go build -o /golangorg -tags=golangorg golang.org/x/website/cmd/golangorg

# Clean up goroot for the final image.
RUN cd /goroot && git clean -xdf

# Add build metadata.
RUN cd /goroot && echo "go repo HEAD: $(git rev-parse HEAD)" >> /goroot/buildinfo
RUN echo "requested go ref: ${GO_REF}" >> /goroot/buildinfo
ARG WEBSITE_HEAD
RUN echo "x/website HEAD: ${WEBSITE_HEAD}" >> /goroot/buildinfo
ARG WEBSITE_CLEAN
RUN echo "x/website clean: ${WEBSITE_CLEAN}" >> /goroot/buildinfo
ARG DOCKER_TAG
RUN echo "image: ${DOCKER_TAG}" >> /goroot/buildinfo
ARG BUILD_ENV
RUN echo "build env: ${BUILD_ENV}" >> /goroot/buildinfo

RUN rm -rf /goroot/.git

# Final image
#############

FROM gcr.io/distroless/base

WORKDIR /app
COPY --from=build /golangorg /app/
COPY --from=build /website/cmd/golangorg/hg-git-mapping.bin /app/

COPY --from=build /goroot /goroot
ENV GOROOT /goroot

COPY --from=build /website/cmd/golangorg/index.split.* /app/
ENV GOLANGORG_INDEX_GLOB index.split.*

CMD ["/app/golangorg"]
