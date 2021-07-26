# Build environment
# -----------------
FROM golang:1.16-alpine as build
LABEL stage=builder

WORKDIR /src

COPY . .

RUN apk add --no-cache git ca-certificates && \
    go env -w GO111MODULE=on && \
    git config --global url."https://{{GITHUB_TOKEN}}:x-oauth-basic@{{REPO_HOST}}/{{REPO_ROOT}}/".insteadOf "https://{{REPO_HOST}}/{{REPO_ROOT}}/" && \
    go env -w GOPRIVATE={{REPO_HOST}}/{{REPO_ROOT}} && \
    go mod tidy && \
    CGO_ENABLED=0 go build -ldflags '-w -s' -o /bin/app

# Deployment environment
# ----------------------
FROM scratch

COPY --from=build /bin/app /bin/app

# Metadata
LABEL org.label-schema.build-date={{TIMESTAMP}} \
      org.label-schema.name={{REPO_NAME}} \
      org.label-schema.description="{{REPO_DESCRIPTION}}" \
      org.label-schema.vcs-url={{REPO_URL}} \
      org.label-schema.vcs-ref={{REPO_COMMIT}} \
      org.label-schema.vendor={{REPO_ROOT}} \
      org.label-schema.version={{REPO_TAG}} \
      org.label-schema.docker.schema-version="1.0"

ARG WEBHOOK_PORT
EXPOSE ${WEBHOOK_PORT}

CMD ["/bin/app"]