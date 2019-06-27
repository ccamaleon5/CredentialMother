# Credential Mother

This is a Provider Credential Server that validates, signs, generates, revokes and updates credential to identify persons, institutions and objects.

The Provider Credential Server sign a credential using its own keys, it is configurable.

The Provider Credential needs to manage its own repository of credentials, default is smart contract that is deployed when server init.

The Credentials are verifiable against blockchain default, but you can configure and choose your proof and revocation list service.

## Prerequisites

* Go 1.12+ installation or later
* **GOPATH** environment variable is set correctly
* docker version 17.03 or later

## Package overview

1. **cmd/credential-provider-server** contains the main for the credential-provider-server command.
2. **lib** contains most of the code.
3. **blockchain** contains smart contracts, ABIs, connections to Ethereum
4. **business** contains business logic that will be consume by APIs
5. **models** conatins data models of requests and responses of APIs
6. **swagger** contains documentation about APIs in Swagger and SwaggerUI to visualize this documentation
7. **util** contains util functions about files and ethereum address

## Install

```
$ git clone https://github.com/ccamaleon5/CredentialMother

$ export GO111MODULE=on

$ cd CredentialMother
$ go build
```

## Run

```
$ credential-provider-server init [-x PASSWORD]
[PASSWORD] is your keystore password that will be created
$ credential-provider-server start --port=8000 --tlscertificate server.crt --tlskey server.key [-x PASSWORD]
```

where --port is a listen port http

You can try in localhost:8000/swagger-ui/

### Docker

* Clone this repository

```
$ git clone https://github.com/ccamaleon5/CredentialProvider
```

* Create a local directory that saves application data  

```
$ mkdir /CredentialData
```

* Copy the configuration file and swaggerui from repository to your local directory created above:

```
$ cp repo/CredentialProvider/credential-provider-server-config-yaml /CredentialData/
$ cp -r repo/CredentialProvider/swagger/swaggerui  /CredentialData/ 

```

* Now pull the docker image and run the container, setting your node identity and the folder location that will be the volume 

```
$ docker pull aparejaa/credentialprovider:1.0.0
$ docker run -dit -v {CredentialProvider_DIR}:/CredentialProvider -p 8000:8000 -p 8001:8001 aparejaa/credentialprovider:1.0.0 credential-provider-server init [-x PASSWORD]
$ docker run -dit -v {CredentialProvider_DIR}:/CredentialProvider -p 8000:8000 -p 8001:8001 aparejaa/credentialprovider:1.0.0
```

* The container will create KeyStore in your local volume

You can try in localhost:8000/swagger-ui/