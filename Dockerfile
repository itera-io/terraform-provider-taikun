FROM golang:1.17.2-alpine AS build

WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/terraform-provider-taikun .

FROM scratch AS bin
COPY --from=build /out/terraform-provider-taikun /
