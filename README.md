# Go Health Checker
> A simple health checker for your web applications.
> 
> Go Health Checker checks your website specified in CLI args every couple of seconds. 
> This is repeated until CTRL+C is pressed. 
> Once that is done, the program will print stats about analysis.

## Run the app
...

## Architecture

The application is built using these components:
 - View - representation of the data to user
   - CLIView - command line interface
 - Model - data structure
   - HealthCheckResult - data structure for health check result
 - Service - business logic
   - HealthCheckService - service for health check, pinging the websites
 - Store - data access layer
   - InMemoryStore - in-memory store for health check results
 - Controller - putting it all together
   - HealthCheckController - controller for health check, responsible for starting and stopping the health check

## How to run the app

```bash
go run cmd/app/main.go \
  https://www.seznam.cz \
  https://www.google.com \
  https://www.cdn77.com \
  https://www.nonexistingdomain.com \
  https://www.youtube.com
```

## What can be improved?

- UI rendering
- E2E testing of the CLI -> run the main.go with some arguments, 
not the tests mainly test the integration of yeah component, but not the CLI itself
- Table rendering - proper info like - TIMEOUT could be displayed
- Better models -> SuccessFull model, Failed model, Timeout model etc, with proper inheritance


## Run the tests

```bash
go test -v ./...
```