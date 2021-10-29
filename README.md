# ledger

## Setup

- install go
- install docker
- install docker-compose
- cd dev-infrastructure
- docker-compose up
- in the project root: `tern migrate --migrations db-migrations`
- start app: `go run .`

## Execute requests

Example requests are present in the `tests.http` file. You can execute execture the reqeuests by using [this](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) vscode plugin or intellij.

## How to scale

Create spezialized ledgers for each sub problem
Post to general ledger after x records or time

https://www.accountingtools.com/articles/how-to-post-to-the-general-ledger.html

## Run migrations

tern migrate --migrations db-migrations
