# Gemini Agent: Project Context - Money Backend

This document captures the unique context of the Money Backend project, providing the Gemini agent with the necessary information to operate effectively within this codebase.

## Project Description
The Money Backend is a personal finance management system built as a set of serverless microservices. It handles core financial entities such as expenses, income, users, periods, and savings goals.

## Architecture Overview
- **Serverless:** The project is primarily built using AWS Lambda functions.
- **Microservices-oriented:** Functions are grouped by domain (e.g., expenses, income, users) within the `api/functions/` directory.
- **Clean Architecture Principles:** The project separates concerns into:
    - `models`: Core domain entities.
    - `usecases`: Business logic and orchestration.
    - `storage`: Data access layer (repositories) with implementations for DynamoDB.
    - `api/functions`: Entry points (Lambda handlers).
    - `shared`: Common utilities, logging, and infrastructure code.
- **Data Persistence:** Uses Amazon DynamoDB as the primary data store.
- **Caching:** Uses Redis for caching (found in `storage/cache`).
- **Observability:** Integrated with the ELK stack for logging and CloudWatch.

## Key Technologies
- **Language:** Go
- **Cloud Provider:** AWS (Lambda, DynamoDB, SQS, CloudWatch, Secrets Manager)
- **Database:** Amazon DynamoDB
- **Caching:** Redis
- **Logging:** ELK Stack (Logstash, Filebeat), CloudWatch
- **Testing:** `testing` package with `testify` (assert/require)

## Key Files and Directories
- `api/functions/`: Contains the Lambda function handlers for various domains.
- `models/`: Defines the data structures for the domain.
- `usecases/`: Contains the core business logic.
- `storage/`: Repository implementations, primarily for DynamoDB.
- `shared/`: Utility packages (logger, router, env, etc.).
- `scripts/`: Deployment scripts for different services.
- `auth/`: Authentication logic, including a custom authenticator and Lambda authorizer.

## Local Development & Setup
- The project uses environment variables managed via a `.env` file in the root.
- Deployment is handled via shell scripts in the `scripts/` directory.
- Build artifacts are stored in `bin/`.

## Project Conventions

### Unit Test Generation
When generating unit tests for this project, Gemini must adhere to the following standards:
- **Assertion Library:** Use the `github.com/stretchr/testify/assert` package.
- **Subtests:** Use `t.Run(test_case_name, func(t *testing.T) { ... })` to structure multiple test cases within a single test function.
- **Mocking:** Utilize existing mocks in `storage/` and `usecases/` or generate new ones using the project's established patterns (e.g., `logger_mock.go`, `dynamodb_mock.go`).
- **Style:** Prefer table-driven tests when testing multiple inputs/outputs for the same logic.

### Path Construction
All file paths used in tool calls must be absolute.
