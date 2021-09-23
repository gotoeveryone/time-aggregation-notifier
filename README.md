# Time aggregation notifier

![Build Status](https://github.com/gotoeveryone/time-aggregation-notifier/workflows/Build/badge.svg)

## Requirements

- Python 3.8
- pipenv
- AWS account (use to DynamoDB, Lambda and Systems Manager)
- Slack account

## Setup

```console
$ pipenv install # When with dev-package add `-d` option.
$ cp .env.example .env # Please edit the value.
```

## Run (Local)

```console
$ pipenv run execute
```

## Code check and format (with pycodestyle and autopep8)

```console
$ # Code check
$ pipenv run code_check
$ # Format
$ pipenv run code_format
```

## Deploy

```console
$ cp .chalice/config.json.example .chalice/config.json # Please edit the value.
$ pipenv run deploy
```
