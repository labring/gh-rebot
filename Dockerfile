FROM alpine:3.16.2

LABEL version="1.0.0"
LABEL repository="https://github.com/labring-actions/gh-rebot"
LABEL homepage="https://github.com/labring-actions/gh-rebot"
LABEL maintainer="Sealos Inc."
LABEL "com.github.actions.name"="Automatic github bot for Sealos"
LABEL "com.github.actions.description"="Automatically github bot for Sealos using comment"

RUN apk --no-cache add jq bash curl git git-lfs github-cli

COPY gh-rebot /usr/bin/

ENTRYPOINT ["/usr/bin/gh-rebot"]
