FROM golang:1.22.2-alpine as builder
LABEL maintainer="Deploy Buddy"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o main ./cmd/app/main.go

# ---

FROM alpine:3.15 as runtime
WORKDIR /app

COPY --from=builder /app/main .
COPY .env .env
COPY metadatas/package.xml metadatas/package.xml

RUN apk add --no-cache nodejs npm && \
    npm install -g sfdx-cli && \
    npm cache clean --force && \
    apk del npm
RUN sfdx --version

# Runtime stage
FROM debian:buster-slim

WORKDIR /app

# Install Salesforce CLI
RUN apt-get update && apt-get install -y curl

RUN curl -sL https://developer.salesforce.com/media/salesforce-cli/sf/channels/stable/sf-linux-x64.tar.gz | tar xz -C /usr/local/bin

ENV PATH=$PATH:/usr/local/bin/sfdx/bin

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]