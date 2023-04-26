ARG TARGETARCH=amd64

FROM --platform=linux/${TARGETARCH} alpine:3.16.2

MAINTAINER "Sealos Inc."

RUN apk --no-cache add jq bash curl git git-lfs github-cli

COPY gh-rebot /usr/bin/

ENTRYPOINT ["/usr/bin/gh-rebot"]

CMD ["--help"]
