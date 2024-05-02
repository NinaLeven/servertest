# Servertest

This repo is a dummy for plaing around with servers.

## Quickstart

```bash
go build -o servertest servertest/bin && ./servertest -c config/dev.yaml
```

## API

```bash
go generate ./...
```

## Docker

```bash
docker build -t servertest .
```

To tag and push the image to the registry:

```bash
VERSION=latest
docker tag servertest:$VERSION vpbukhti/servertest:$VERSION
docker push vpbukhti/servertest:$VERSION
```