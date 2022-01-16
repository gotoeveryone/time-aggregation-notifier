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

Use [lambroll](https://github.com/fujiwara/lambroll).

```console
$ cp deploy/function.json.example deploy/function.json # Please edit the value.
$ go build -o deploy/time-aggregation-notifier ./src/main.go
$ cd deploy
$ lambroll deploy
```
