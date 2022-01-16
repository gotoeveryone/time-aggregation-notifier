# Time aggregation notifier

![Build Status](https://github.com/gotoeveryone/time-aggregation-notifier/workflows/Build/badge.svg)

## Requirements

- Golang
- AWS account (use to DynamoDB, Lambda and Systems Manager)
- Slack account

## Setup

```console
$ go mod download
$ cp .env.example .env # Please edit the value.
```

## Run (Local)

```console
$ go run src/main.go
```

## Deploy

```console
$ cp .chalice/config.json.example .chalice/config.json # Please edit the value.
$ pipenv run deploy
```
