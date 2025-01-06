# Order Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles product operations like create, read, update, and delete. Using AWS DynamoDB as main database to storing data like name, categories, stock, etc. Also using redis as cache database for categories for low latency and high performance purpose.

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/
├── cmd/
│   └── app/
├── config/
├── internal/   
│   ├── app/
│   ├── constant/
│   ├── controller/
│   │   ├── http/
│   │   |   └── v1/
│   │   └── kafka
│   │       └── v1/
│   ├── dto/
│   ├── entity/
│   ├── usecase/
│   │   └── repo/
│   ├── usecase/
│   └── utils/
├── migrations/
└── pkg/
    ├── aws/
    ├── httpserver/
    ├── kafka/
    ├── logger/
    └── redis/
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