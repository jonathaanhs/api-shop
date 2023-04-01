# Shop

Checkout services for shopping with promos.

## Requirements
- Go 1.17
- [golang-migrate](https://github.com/golang-migrate/migrate)
- Docker
- Docker Compose

## Quick Start
```bash
make test # to run unit test
make run-pg # run the postgres db
make stop-pg # stop run the postgres db
make run # run the program
```

## Call API
```bash
Case 1: Buying more than 3 Alexa Speakers will have a 10% discount on all Alexa speakers

curl --location 'http://localhost:8089/graphql' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n\tcheckout(items: [{product_id: 3, qty: 3}]) {\n\t\titems\n\t\ttotal_amount\n\t}\n}","variables":{}}'

{
    "data": {
        "checkout": {
            "items": [
                "Alexa Speaker",
                "Alexa Speaker",
                "Alexa Speaker"
            ],
            "total_amount": 295.65
        }
    }
}
```

```bash
Case 2: Buy 3 Google Homes for the price of 2

curl --location 'http://localhost:8089/graphql' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n\tcheckout(items: [{product_id: 1, qty: 3}]) {\n\t\titems\n\t\ttotal_amount\n\t}\n}","variables":{}}'

{
    "data": {
        "checkout": {
            "items": [
                "Google Home",
                "Google Home",
                "Google Home"
            ],
            "total_amount": 99.98
        }
    }
}
```

```bash
Case 3: Each sale of a MacBook Pro comes with a free Raspberry Pi B

curl --location 'http://localhost:8089/graphql' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n\tcheckout(items: [{product_id: 2, qty: 1}]) {\n\t\titems\n\t\ttotal_amount\n\t}\n}","variables":{}}'

{
    "data": {
        "checkout": {
            "items": [
                "MacBook Pro",
                "Raspberry Pi B"
            ],
            "total_amount": 5399.99
        }
    }
}
```

```bash
Case 4: the product qty is not enough to fulfill the request

curl --location 'http://localhost:8089/graphql' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n\tcheckout(items: [{product_id: 2, qty: 100}]) {\n\t\titems\n\t\ttotal_amount\n\t}\n}","variables":{}}'

{
    "data": {
        "checkout": null
    },
    "errors": [
        {
            "message": "the product MacBook Pro qty is not enough to fulfill the request",
            "locations": [
                {
                    "line": 2,
                    "column": 2
                }
            ],
            "path": [
                "checkout"
            ]
        }
    ]
}
```

## Unit Test Coverage

```console
foo@bar:~$ make test     
go clean -testcache
go test ./... --cover
?       github.com/learn/api-shop/cmd   [no test files]
?       github.com/learn/api-shop/internal      [no test files]
ok      github.com/learn/api-shop/internal/controller   1.178s  coverage: 85.7% of statements
?       github.com/learn/api-shop/internal/generated/mock       [no test files]
?       github.com/learn/api-shop/internal/infra        [no test files]
ok      github.com/learn/api-shop/internal/repo 0.684s  coverage: 100.0% of statements
ok      github.com/learn/api-shop/internal/service      0.990s  coverage: 100.0% of statements
ok      github.com/learn/api-shop/pkg/sqlkit    0.477s  coverage: 100.0% of statements
```