# Order Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles product operations like create, read, update, and delete. Using AWS DynamoDB as main database to storing data like name, categories, stock, etc. Also using redis as cache database for categories for low latency and high performance purpose.

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/      # github workflows to automatically test, build, and push
├── cmd/
│   └── app/            # configuration and log initialization
├── config/             # configuration
├── internal/   
│   ├── app/            # one run function in the `app.go`
│   ├── constant/       # global constant
│   │   ├── http/
│   │   |   └── v1/     # rest http
│   │   └── kafka
│   │       └── v1/     # kafka consumer
│   ├── entity/         # entities of business logic (models) can be used in any layer
│   ├── usecase/        # business logic
│   │   └── repo/       # abstract storage (database) that business logic works with
│   └── utils/          # helpers function
└── pkg/
    ├── aws/            # aws initialization for client, dynamodb, and s3
    ├── httpserver/     # http server initialization
    ├── kafka/          # kafka initialization
    ├── logger/         # logger initialization
    └── redis/          # redis initialization
```

## Tech Stack
- Backend: Go
- Authorization: AWS Cognito
- Database: AWS DynamoDB and Redis
- CI/CD: Github Actions
- Message Broker: Apache Kafka
- Container: Docker

## API Documentation
tbd