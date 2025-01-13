# Dating App Backend

A backend system for a dating app, built with **Gin** and designed using the **Clean Architecture** pattern.

## Features

- User sign-up and login with JWT authentication.
- Daily swipe quota for free users.
- Premium packages to unlock features like unlimited swipes.
- Modular and scalable architecture.

---

## Directory Structure

```plaintext
dating-app-backend/
├── .vscode                 # Debug mode configuration
├── api/                    # Define the service want to run
|   └── http                # To run http server
├── cmd/                    # Entry point of the application
│   └── main.go
├── configs/                # Configuration files
├── docker                  # Docker configuration
├── internal/               # Application business logic
│   ├── constants/          # The constant variables
│   ├── domain/             # Entities and interfaces
│   ├── dto/                # Data validation
│   ├── usecase/            # Business rules and interactions
│   ├── delivery/http/      # HTTP handlers (Gin routes)
│   ├── repository/         # Database operations
├── pkg/                    # Shared utilities
│   ├── converter/          # Converter utils
│   ├── database/
│       ├── postgres        # Utils for easier database transaction
│   ├── generator/          # Generator utils
│   ├── jwt/                # JWT generation and validation
│   └── logger/             # Logger setup with new relic
│   ├── object_storage/     # Object storage utils, to upload the image, for now using imagekit
├── go.mod                  # Go modules file
├── Makefile                # Easier to run the service
└── README.md               # Documentation
```

## How to run

### HTTP Server

1. Make sure you've already fill .env file.
2. Run `make run-http` to start the http service.

### Test

#### Unit test

1. Run `make test` to start unit test. For now, it's small coverage because of the short time of development.
2. Run `make coverage` to see detail about the unit test result coverage.

#### Postman Test

You can test through postman with this [link](https://documenter.getpostman.com/view/3884681/2sAYQWLE9m).

### Integration Test

You can test with Postman to test the integration or integrate with Frontend or Mobile.

## Additional Enhancements

- **Linting**: Added `golangci-lint` to ensure code quality and consistency.
- **Deployment**: Deployment with Railway app. For the future, it can be deployed as kubernetes pod.
- **Cache**: Use redis for some case. Still small using because of the short time of development.
- **Logging**: Use new relic for logging.
- **Object Storage**: Use object storage to save the image.
