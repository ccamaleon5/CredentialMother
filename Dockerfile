# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git \
    && git clone https://github.com/ccamaleon5/CredentialMother.git
WORKDIR /go/CredentialMother/cmd/credential-provider-server
RUN export GO111MODULE=on && go build

# Pull Credential Provider Server into a second stage deploy alpine container
FROM alpine:3.9
RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

RUN mkdir -p /CredentialMother 
WORKDIR /CredentialMother

COPY --from=builder go/CredentialMother/cmd/credential-provider-server /usr/local/bin/
COPY --from=builder go/CredentialMother/cmd/credential-provider-server/server.crt /CredentialMother/
COPY --from=builder go/CredentialMother/cmd/credential-provider-server/server.key /CredentialMother/
COPY --from=builder go/CredentialMother/swagger/ /CredentialMother/

EXPOSE 8000 8001
CMD ["sh","-c","credential-provider-server start --host=0.0.0.0 --tlscertificate=/usr/local/bin/server.crt --tlskey=/usr/local/bin/server.key --port 8000 --secret Password."]